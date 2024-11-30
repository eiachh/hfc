package service

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/eiachh/hfc/logger"
	"github.com/eiachh/hfc/types"
)

type ChatGptAiCaller struct {
	api_key        string
	url            string
	callCountLimit int
}

func NewChatGptAiCaller() *ChatGptAiCaller {
	return &ChatGptAiCaller{
		callCountLimit: 3,
		api_key:        getApiKey(),
		url:            "https://api.openai.com/v1/chat/completions",
	}
}

func (ai *ChatGptAiCaller) ParseOff(trimmedOffByte []byte) (*types.Product, error) {
	var offProd types.Product
	json.Unmarshal(trimmedOffByte, &offProd)
	barC := offProd.Code

	messages := []types.AiMessage{
		DefaultParseTaskAiMessage(),
		{Role: "user", Content: string(trimmedOffByte)},
	}
	aibody := NewAiReqBody("gpt-4o", messages)
	aibodyJson, _ := json.MarshalIndent(aibody, "", "  ")

	chatComp, err := ai.callGpt(aibodyJson)
	if err != nil {
		return nil, err
	}

	return ai.parseAiResp(barC, chatComp, aibody, 1)
}

func (ai *ChatGptAiCaller) WebScrapeParse(barcode int64, chatComp *types.ChatCompletion, aibody *types.AiReqBody, callCount int) (*types.Product, error) {
	logger.Log().Debug("Running webscraping")
	allowedWebscrapeLimit := ai.callCountLimit - callCount

	if aibody != nil && chatComp == nil {
		logger.Log().Fatalf("Invalid state aibody not nil but comp is. existingAiBody: %s", aibody)
	} else if aibody == nil && chatComp != nil {
		logger.Log().Fatalf("Invalid state aibody is nil but comp is NOT. ChatCompletion: %s", chatComp)
	}

	scrapedHtmlText := ScrapeDataOf(barcode, callCount)
	logger.Log().Debugf("Got back scraped html: %s", scrapedHtmlText)

	if aibody == nil && chatComp == nil {
		messages := []types.AiMessage{
			DefaultParseTaskAiMessage(),
			{Role: "user", Content: scrapedHtmlText},
			{Role: "system", Content: ("Remember you are filling out a product info of a food or drink product, ignore the text unrelated to any food or drink item. You got some scraped data, try to logically fill out the final_response. Mock, empty, placeholder data is still not allowed. Try to figure out the fields logically. If you do believe the scraped data is not about a the food product, you can still use request_more_info, or if you think you still need more info. If you can confidently answer it is prefered but you can still use the request_more_info function " + strconv.Itoa(allowedWebscrapeLimit) + " more times.")},
		}
		aibody = NewAiReqBody("gpt-4o-mini", messages)
	} else {
		ai.continueAiBody(chatComp, aibody, callCount, scrapedHtmlText)
	}

	aiBodyJson, _ := json.MarshalIndent(aibody, "", "  ")
	respChatComp, err := ai.callGpt(aiBodyJson)
	logger.Log().Debugf("Got response: %s", respChatComp)

	if err != nil {
		return nil, err
	}
	return ai.parseAiResp(barcode, respChatComp, aibody, 1)
}

func (ai *ChatGptAiCaller) callGpt(reqBodyJson []byte) (*types.ChatCompletion, error) {
	logger.Log().Debugf("Calling openapi with request: %s", reqBodyJson)
	client := &http.Client{}

	req, err := http.NewRequest("POST", ai.url, bytes.NewBuffer(reqBodyJson))
	if err != nil {
		return &types.ChatCompletion{}, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ai.api_key)

	resp, err := client.Do(req)
	if err != nil {
		return &types.ChatCompletion{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &types.ChatCompletion{}, err
	}

	logger.Log().Debugf("Openapi response: %s", body)
	var chatComp types.ChatCompletion
	json.Unmarshal(body, &chatComp)
	return &chatComp, nil
}

func (ai *ChatGptAiCaller) continueAiBody(chatComp *types.ChatCompletion, aiBody *types.AiReqBody, callCount int, scrapedHtmlText string) *types.AiReqBody {
	allowedWebscrapeLimit := ai.callCountLimit - callCount
	var sysInstructionAfterToolCall string
	if callCount == ai.callCountLimit {
		sysInstructionAfterToolCall = "You got some additional scraped data from the function, try to logically fill out the final_response. Mock, empty, placeholder data is still not allowed. Try to figure out the fields logically. Use final_response function to fill out the data, you can fill out the fields based on assumption, you are not allowed to fill the brand if you would just guess, write UNKNOWN to fields that would be just guesses without any bases. You are not allowed to use request_more_info anymore!"
	} else {
		sysInstructionAfterToolCall = ("Remember you are filling out a product info of a food or drink product, ignore the text unrelated to any food or drink item. You got some scraped data, try to logically fill out the final_response. Mock, empty, placeholder data is still not allowed. Try to figure out the fields logically. If you do believe the scraped data is not about a the food product, you can still use request_more_info, or if you think you still need more info. If you can confidently answer it is prefered but you can still use the request_more_info function " + strconv.Itoa(allowedWebscrapeLimit) + " more times.")
	}

	sysMsg := types.AiMessage{
		Role:    "system",
		Content: sysInstructionAfterToolCall,
	}
	scrapeMsg := types.AiMessage{
		Role:       "tool",
		Content:    scrapedHtmlText,
		Name:       "request_more_info",
		ToolCallId: chatComp.Choices[0].Message.ToolCalls[0].ID,
	}

	aiBody.Messages = append(aiBody.Messages, chatComp.Choices[0].Message)
	aiBody.Messages = append(aiBody.Messages, scrapeMsg)
	aiBody.Messages = append(aiBody.Messages, sysMsg)
	return aiBody
}

func getApiKey() string {
	return os.Getenv("OPENAI_API_KEY")
}

func (ai *ChatGptAiCaller) parseAiResp(barC int64, chatComp *types.ChatCompletion, aibody *types.AiReqBody, callCount int) (*types.Product, error) {
	var prod types.Product

	// TODO Handle logging if by any chance multiple choice or tool calls happened.
	if len(chatComp.Choices) > 1 {
		logger.Log().Fatal("Multiple choices were given by the AI")
	}
	if len(chatComp.Choices[0].Message.ToolCalls) > 1 {
		logger.Log().Fatal("Multiple toolCalls were given by the AI")
	}

	if chatComp.Choices[0].Message.ToolCalls[0].Function.Name == "request_more_info" {
		return ai.WebScrapeParse(barC, chatComp, aibody, callCount)
	} else {
		if responseArgErr := json.Unmarshal([]byte(chatComp.Choices[0].Message.ToolCalls[0].Function.Arguments), &prod); responseArgErr != nil {
			logger.Log().Error(responseArgErr)
			return nil, responseArgErr
		}
	}
	prod.Code = barC
	if prod.Brands != "UNKNOWN" {
		prod.DisplayName = prod.Brands + " " + prod.Name
	} else {
		prod.DisplayName = prod.Name
	}
	return &prod, nil
}
