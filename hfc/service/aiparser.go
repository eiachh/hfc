package service

import (
	"encoding/json"

	"github.com/eiachh/hfc/types"
	"github.com/labstack/gommon/log"
)

const (
	jsonDataMOCK = `{
	"model": "gpt-4o",
	"max_completion_tokens": 250,
	"frequency_penalty": 0.0,
	"top_p": 0.1,
	"messages": [
		{
			"role": "system",
			"content": "Your task is to parse the data about food and drink products from the user into a json fill out the final_response function when you think you have enough data to do so, for categories_hierarchy if the given infos are enough try to guess it, you can also request more data if the given input was not enough, only use this in last resort. Mock data empty values and placeholders are strongly penalized. Always use english if needed translate. Be pessimistic with the expire_avg guess lower end is prefered."
		},
		{
			"role": "user",
			"content": " code:4014500513010, _keywords: [dairy,dessert,fermented,food,jogobella,jogurt,jogurty,lesní,milk,mléčné,ovoce,ovocné,product,výrobky,zott], brands: Zott, product_name: Jogobella lesní ovoce"
		}
	],
	"tools": [
		{
			"type": "function",
			"function": {
				"name": "final_response",
				"description": "After you have enough data to fill out all the required fields in this function do it and you succeed with the task.",
				"parameters": {
					"type": "object",
					"properties": {
						"brand": {
							"type": "string",
							"description": "The brand of the food or drink product"
						},
						"product_name": {
							"type": "string",
							"description": "The english name of the product make it short but easy to understand."
						},
						"categories_hierarchy": {
							"description": "Big to small food groups, make sure the groups are ordered. You can fill this logically if other info is enough. Example for a fruit yoghurt: categories_hierarchy: [dairies,fermented-foods,fermented-milk-products,desserts,dairy-desserts,fermented-dairy-desserts,fermented-dairy-desserts-with-fruits,yogurts,fruit-yogurts]",
							"type": "array",
							"items": {
								"type": "string"
							}
						},
						"expire_avg": {
							"type": "string",
							"description": "Assuming the product was made today the expected time in days until it expires only the numbers."
						},
						"measurement_unit": {
							"type": "string",
							"description": "The most logical unit that this product should be measured in, during a day to day conversation. Try to guess the most fitting one.",
							"enum": [
								"milliliter",
								"gramm",
								"piece"
							]
						}
					},
					"required": [
						"brand",
						"product_name",
						"categories_hierarchy",
						"expire_avg",
						"measurement_unit"
					]
				}
			}
		},
		{
			"type": "function",
			"function": {
				"name": "request_more_info",
				"description": "Request more info if the already given is not enough to call final_response",
				"parameters": {
					"type": "object",
					"properties": {
						"request": {
							"type": "string",
							"description": "Fill it with true if you need more data."
						}
					},
					"required": [
						"request"
					]
				}
			}
		}
	],
	"tool_choice": "auto"
}`
)

type AiParser struct {
	caller AiCaller
}

func NewAiParser(ai AiCaller) *AiParser {
	return &AiParser{
		caller: ai,
	}
}

func (ai *AiParser) ConvertOffToLocCache(offByte *[]byte) *types.Product {
	trimmedOff, err := trimOff(*offByte)
	if err != nil {
		log.Error(err)
	}
	ai.letAiFill(&trimmedOff)
	return nil
}

func trimOff(input []byte) ([]byte, error) {
	keysToKeep := []string{"code", "_keywords", "brands", "categories_hierarchy", "product_name"}
	var data map[string]interface{}

	if err := json.Unmarshal(input, &data); err != nil {
		return nil, err
	}

	filteredData := make(map[string]interface{})
	for _, key := range keysToKeep {
		if value, exists := data[key]; exists {
			filteredData[key] = value
		}
	}

	result, err := json.Marshal(filteredData)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ai *AiParser) letAiFill(trimmedOffByte *[]byte) *types.Product {
	ai.caller.callAI(*trimmedOffByte)
	return nil
}
