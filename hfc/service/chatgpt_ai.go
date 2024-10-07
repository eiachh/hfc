package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/eiachh/hfc/types"
	"github.com/labstack/gommon/log"
)

type aiReqBody struct {
	Model               string      `json:"model"`
	MaxCompletionTokens int         `json:"max_completion_tokens"`
	FrequencyPenalty    float32     `json:"frequency_penalty"`
	TopP                float32     `json:"top_p"`
	Messages            []aiMessage `json:"messages"`
	Tools               []aiTool    `json:"tools"`
}

type aiMessage struct {
	Role       string     `json:"role"`
	Content    string     `json:"content,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls"`
	Refusal    *string    `json:"refusal,omitempty"`
	Name       string     `json:"name,omitempty"`
	ToolCallId string     `json:"tool_call_id,omitempty"`
}

type aiTool struct {
	Type     string     `json:"type"`
	Function aiFunction `json:"function"`
}

type aiFunction struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Parameters  aiParameters `json:"parameters"`
}

type aiParameters struct {
	Type       string                `json:"type"`
	Properties map[string]aiProperty `json:"properties"`
	Required   []string              `json:"required"`
}

type aiProperty struct {
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Enum        []string    `json:"enum,omitempty"`
	Items       *aiProperty `json:"items,omitempty"`
}

// RESPONSE PARSER
type ChatCompletion struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
	SystemFingerprint string   `json:"system_fingerprint"`
}

type Choice struct {
	Index        int       `json:"index"`
	Message      aiMessage `json:"message"`
	Logprobs     *string   `json:"logprobs,omitempty"`
	FinishReason string    `json:"finish_reason"`
}

type ToolCall struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

type Function struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type Usage struct {
	PromptTokens            int                     `json:"prompt_tokens"`
	CompletionTokens        int                     `json:"completion_tokens"`
	TotalTokens             int                     `json:"total_tokens"`
	PromptTokensDetails     PromptTokensDetails     `json:"prompt_tokens_details"`
	CompletionTokensDetails CompletionTokensDetails `json:"completion_tokens_details"`
}

type PromptTokensDetails struct {
	CachedTokens int `json:"cached_tokens"`
}

type CompletionTokensDetails struct {
	ReasoningTokens int `json:"reasoning_tokens"`
}

func NewAiReqBody() aiReqBody {
	return aiReqBody{
		Model:               "",
		MaxCompletionTokens: 250,
		FrequencyPenalty:    0.0,
		TopP:                0.1,
		Messages: []aiMessage{
			{
				Role:    "system",
				Content: "Your task is to parse the data about food and drink products from the user into a json fill out the final_response function when you think you have enough data to do so, for categories_hierarchy if the given infos are enough try to guess it, you are not allowed to use message content only use the functions. If the given input was not enough use the function request_more_info, only use this if you cannot figure out logically the needed fields. Mock data empty values and placeholders are strongly penalized. Always use english if needed translate. Be pessimistic with the expire_avg guess lower end is prefered.",
			},
			{
				Role:    "user",
				Content: " code:4014500513010, _keywords: [dairy,dessert,fermented,food,jogobella,jogurt,jogurty,lesní,milk,mléčné,ovoce,ovocné,product,výrobky,zott], brands: Zott, product_name: Jogobella lesní ovoce",
			},
		},
		Tools: []aiTool{
			{
				Type: "function",
				Function: aiFunction{
					Name:        "final_response",
					Description: "After you have enough data to fill out all the required fields in this function do it and you succeed with the task.",
					Parameters: aiParameters{
						Type: "object",
						Properties: map[string]aiProperty{
							"brands": {
								Type:        "string",
								Description: "The brand of the food or drink product",
							},
							"product_name": {
								Type:        "string",
								Description: "The english name of the product make it short but easy to understand.",
							},
							"categories_hierarchy": {
								Type:        "array",
								Description: "Big to small food groups, make sure the groups are ordered. You can fill this logically if other info is enough. Example for a fruit yoghurt: categories_hierarchy: [dairies,fermented-foods,fermented-milk-products,desserts,dairy-desserts,fermented-dairy-desserts,fermented-dairy-desserts-with-fruits,yogurts,fruit-yogurts]",
								Items:       &aiProperty{Type: "string"},
							},
							"expire_avg": {
								Type:        "string",
								Description: "Assuming the product was made today the expected time in days until it expires only the numbers.",
							},
							"measurement_unit": {
								Type:        "string",
								Description: "The most logical unit that this product should be measured in, during a day to day conversation. Try to guess the most fitting one.",
								Enum:        []string{"milliliter", "gramm", "piece"},
							},
						},
						Required: []string{
							"brand",
							"product_name",
							"categories_hierarchy",
							"expire_avg",
							"measurement_unit",
						},
					},
				},
			},
			{
				Type: "function",
				Function: aiFunction{
					Name:        "request_more_info",
					Description: "Request more info if the already given is not enough to call final_response",
					Parameters: aiParameters{
						Type: "object",
						Properties: map[string]aiProperty{
							"request": {
								Type:        "string",
								Description: "Fill it with true if you need more data.",
							},
						},
						Required: []string{"request"},
					},
				},
			},
		},
	}
}

type ChatGptAiCaller struct {
	api_key        string
	url            string
	callCountLimit int
}

func NewChatGptAiCaller() *ChatGptAiCaller {
	return &ChatGptAiCaller{
		callCountLimit: 1,
		api_key:        getApiKey(),
		url:            "https://api.openai.com/v1/chat/completions",
	}
}

func (ai *ChatGptAiCaller) ParseOff(trimmedOffByte []byte) (*types.Product, error) {
	var offProd types.Product
	json.Unmarshal(trimmedOffByte, &offProd)
	barC := offProd.Code

	forceModel := "gpt-4o"
	aibody := NewAiReqBody()
	aibody.Model = forceModel
	// TODO unmock
	aibody.Messages[1].Content = " code:4014500513010"
	aibodyJson, _ := json.MarshalIndent(aibody, "", "  ")

	// Perform the request
	// Read the response body
	//os.WriteFile("chatgptresponse.txt", body, 0644)
	chatComp, err := ai.CallGpt(aibodyJson)
	if err != nil {
		return nil, err
	}

	return ai.parseAiResp(barC, chatComp, &aibody, 1)
}

func (ai *ChatGptAiCaller) WebScrapeParse(barcode int) (*types.Product, error) {
	aibody := NewAiReqBody()
	forceModel := "gpt-4o-mini"
	aibody.Model = forceModel

	rawHtml, scrapeError := ScrapeDataOf(barcode, 1)

	if scrapeError != nil {
		return nil, scrapeError
	}

	firstWebscrapeMsg := aiMessage{
		Role:    "user",
		Content: rawHtml,
	}
	sysMsg := aiMessage{
		Role:    "system",
		Content: "You got some scraped data, try to logically fill out the final_response. Mock, empty, placeholder data is still not allowed. Try to figure out the fields logically. If you do believe the scraped data is not about a the food product, you can still use request_more_info, or if you think you still need more info. the request_more_info function should be used only if you are certain you need more info.",
	}
	aibody.Messages = append(aibody.Messages, firstWebscrapeMsg)
	aibody.Messages = append(aibody.Messages, sysMsg)

	aiBodyJson, _ := json.MarshalIndent(aibody, "", "  ")
	respChatComp, err := ai.CallGpt(aiBodyJson)
	log.Debug(respChatComp)
	if err != nil {
		return nil, err
	}
	log.Debug(respChatComp)
	return ai.parseAiResp(barcode, respChatComp, &aibody, 1)
}

func (ai *ChatGptAiCaller) CallGpt(reqBodyJson []byte) (*ChatCompletion, error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", ai.url, bytes.NewBuffer(reqBodyJson))
	if err != nil {
		return &ChatCompletion{}, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ai.api_key)

	resp, err := client.Do(req)
	if err != nil {
		return &ChatCompletion{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &ChatCompletion{}, err
	}

	log.Info((string(body)))
	var chatComp ChatCompletion
	os.WriteFile("chatgptREQUEST.txt", reqBodyJson, 0644)
	os.WriteFile("chatgptRESPONSE3.txt", body, 0644)
	json.Unmarshal(body, &chatComp)
	return &chatComp, nil
}

// TODO clean, also they cant be nil
func (ai *ChatGptAiCaller) webScrapeWithCtx(barcode int, chatComp *ChatCompletion, aiReqBody *aiReqBody, callCount int) (*types.Product, error) {
	if callCount > ai.callCountLimit {
		return nil, errors.New("tried to call the ai api more than allowed")
	}

	aibody := *aiReqBody
	forceModel := "gpt-4o-mini"
	aibody.Model = forceModel
	aibody.Messages = append(aibody.Messages, chatComp.Choices[0].Message)

	rawHtml, scrapeError := ScrapeDataOf(barcode, callCount)

	sysMsg := aiMessage{
		Role:    "system",
		Content: "The system cannot get more info about the product, try to guess a value logically, if you think you are not confident enough in the guess write 'UNKNOWN'.",
	}
	scrapeMsg := aiMessage{
		Role:       "tool",
		Content:    "Cannot get any more info.",
		Name:       "request_more_info",
		ToolCallId: chatComp.Choices[0].Message.ToolCalls[0].ID,
	}
	if scrapeError == nil {
		scrapeMsg.Content = rawHtml
		sysMsg.Content = "You got some additional scraped data from the function, try to logically fill out the final_response. Mock, empty, placeholder data is still not allowed. Try to figure out the fields logically. If you do believe the scraped data is not about a the food product, you can still use request_more_info, or if you think you still need more info. The request_more_info function can be used if you believe most of the data is still uncertain, otherwise use final_response."
		aibody.Messages = append(aibody.Messages, scrapeMsg)
	}
	if callCount == ai.callCountLimit {
		sysMsg.Content = "You got some additional scraped data from the function, try to logically fill out the final_response. Mock, empty, placeholder data is still not allowed. Try to figure out the fields logically. Use final_response function to fill out the data, you can fill out the fields based on assumption, you are not allowed to fill the brand if you would just guess, write UNKNOWN to fields that would be just guesses without any bases."
	}

	aibody.Messages = append(aibody.Messages, sysMsg)

	aiBodyJson, _ := json.MarshalIndent(aibody, "", "  ")
	respChatComp, err := ai.CallGpt(aiBodyJson)
	if err != nil {
		return nil, err
	}
	log.Debug(respChatComp)
	return ai.parseAiResp(barcode, respChatComp, &aibody, callCount+1)
}

func getApiKey() string {
	return os.Getenv("OPENAI_API_KEY")
}

func (ai *ChatGptAiCaller) parseAiResp(barC int, chatComp *ChatCompletion, aibody *aiReqBody, callCount int) (*types.Product, error) {
	var prod types.Product

	// TODO Handle logging if by any chance multiple choice or tool calls happened.
	if len(chatComp.Choices) > 1 {
		log.Fatal("Multiple choices were given by the AI")
	}
	if len(chatComp.Choices[0].Message.ToolCalls) > 1 {
		log.Fatal("Multiple toolCalls were given by the AI")
	}

	if chatComp.Choices[0].Message.ToolCalls[0].Function.Name == "request_more_info" {
		return ai.webScrapeWithCtx(barC, chatComp, aibody, callCount)
	} else {
		if responseArgErr := json.Unmarshal([]byte(chatComp.Choices[0].Message.ToolCalls[0].Function.Arguments), &prod); responseArgErr != nil {
			log.Error(responseArgErr)
			return nil, responseArgErr
		}
	}
	prod.Code = barC
	if prod.Brands != "UNKNOWS" {
		prod.DisplayName = prod.Brands + " " + prod.Name
	} else {
		prod.DisplayName = prod.Name
	}
	return &prod, nil
}
