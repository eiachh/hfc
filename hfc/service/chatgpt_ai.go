package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
	Role    string `json:"role"`
	Content string `json:"content"`
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

func NewAiReqBody() aiReqBody {
	return aiReqBody{
		Model:               "gpt-4o",
		MaxCompletionTokens: 250,
		FrequencyPenalty:    0.0,
		TopP:                0.1,
		Messages: []aiMessage{
			{
				Role:    "system",
				Content: "Your task is to parse the data about food and drink products from the user into a json fill out the final_response function when you think you have enough data to do so, for categories_hierarchy if the given infos are enough try to guess it, you can also request more data if the given input was not enough, only use this in last resort. Mock data empty values and placeholders are strongly penalized. Always use english if needed translate. Be pessimistic with the expire_avg guess lower end is prefered.",
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
							"brand": {
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
	api_key string
}

func NewChatGptAiCaller() *ChatGptAiCaller {
	return &ChatGptAiCaller{
		api_key: getApiKey(),
	}
}

func (ai *ChatGptAiCaller) callAI(trimmedOffByte []byte) (*[]byte, error) {
	url := "https://api.openai.com/v1/chat/completions"
	client := &http.Client{}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(getAiReqBodyByte()))
	if err != nil {
		return nil, err
	}

	// Set headers if necessary
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ai.api_key)

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	asddsdsd := (string(body))
	fmt.Println(asddsdsd)
	return &body, nil
}

func getApiKey() string {
	return os.Getenv("OPENAI_API_KEY")
}

func getAiReqBodyByte() []byte {
	aibody := NewAiReqBody()
	jsonData, _ := json.MarshalIndent(aibody, "", "  ")
	return jsonData
}
