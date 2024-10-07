package service

import (
	"encoding/json"

	"github.com/eiachh/hfc/storage"
	"github.com/eiachh/hfc/types"
)

type AiParser struct {
	caller AiCaller
	mongo  *storage.MongoStorage
}

func NewAiParser(ai AiCaller, db *storage.MongoStorage) *AiParser {
	return &AiParser{
		caller: ai,
		mongo:  db,
	}
}

func (ai *AiParser) ConvertOffToLocCache(offByte *[]byte) (*types.Product, error) {
	trimmedOff, err := trimOff(*offByte)
	if err != nil {
		return nil, err
	}
	product, err := ai.caller.ParseOff(trimmedOff)
	if err != nil {
		return product, err
	}
	if err := ai.mongo.NewProduct(product); err != nil {
		return product, err
	}
	return product, nil
}

func (ai *AiParser) DoWebscrape(barcode int) (*types.Product, error) {
	return ai.caller.WebScrapeParse(barcode)
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
