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

// Converts offByte ( open food facts db value ) to loc-cache style data, if needed does webscrape for extra info
func (ai *AiParser) ConvertOffToLocCache(offByte *[]byte) (*types.Product, error) {
	trimmedOff, err := trimOff(*offByte)
	if err != nil {
		return nil, err
	}
	return ai.caller.ParseOff(trimmedOff)
}

// Does webscrape from the getgo if off had no info of the product
func (ai *AiParser) DoWebscrape(barcode int64) (*types.Product, error) {
	return ai.caller.WebScrapeParse(barcode, nil, nil, 1)
}

// Only leaves the "useful" data in the off input bytes
func trimOff(input []byte) ([]byte, error) {
	keysToKeep := []string{"code", "_keywords", "brands", "categories_hierarchy", "product_name", "product_quantity", "product_quantity_unit"}
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
