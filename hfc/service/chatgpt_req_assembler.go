package service

import "github.com/eiachh/hfc/types"

func NewAiReqBody(model string, messages []types.AiMessage) *types.AiReqBody {
	return &types.AiReqBody{
		Model:               model,
		MaxCompletionTokens: 250,
		FrequencyPenalty:    0.0,
		TopP:                0.1,
		Messages:            messages,
		Tools: []types.AiTool{
			FinalResponseAiTool(),
			RequestMoreInfoAiTool(),
		},
	}
}

func DefaultParseTaskAiMessage() types.AiMessage {
	return types.AiMessage{
		Role:    "system",
		Content: "Your task is to parse the data about food and drink products from the user into a json fill out the final_response function when you think you have enough data to do so, for categories_hierarchy if the given infos are enough try to guess it, you are not allowed to use message content only use the functions. If the given input was not enough use the function request_more_info, only use this if you cannot figure out logically the needed fields. Mock data empty values and placeholders are strongly penalized. Always use english if needed translate. Be pessimistic with the expire_avg guess lower end is prefered.",
	}
}

func FinalResponseAiTool() types.AiTool {
	return types.AiTool{
		Type: "function",
		Function: types.AiFunction{
			Name:        "final_response",
			Description: "After you have enough data to fill out all the required fields in this function do it and you succeed with the task.",
			Parameters: types.AiParameters{
				Type: "object",
				Properties: map[string]types.AiProperty{
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
						Items:       &types.AiProperty{Type: "string"},
					},
					"expire_avg": {
						Type:        "integer",
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
	}
}

func RequestMoreInfoAiTool() types.AiTool {
	return types.AiTool{
		Type: "function",
		Function: types.AiFunction{
			Name:        "request_more_info",
			Description: "Request more info if the already given is not enough to call final_response",
			Parameters: types.AiParameters{
				Type: "object",
				Properties: map[string]types.AiProperty{
					"request": {
						Type:        "string",
						Description: "Fill it with true if you need more data.",
					},
				},
				Required: []string{"request"},
			},
		},
	}
}
