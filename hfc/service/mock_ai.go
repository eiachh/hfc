package service

import (
	"encoding/json"

	"github.com/eiachh/hfc/types"
	"github.com/labstack/gommon/log"
)

type MockAiCaller struct{}

func (mockAi *MockAiCaller) ParseOff(trimmedOff []byte) (*types.Product, error) {
	//{"code", "_keywords", "brands", "categories_hierarchy", "product_name"}
	var data map[string]interface{}
	retData := make(map[string]interface{})
	if err := json.Unmarshal(trimmedOff, &data); err != nil {
		return nil, err
	}
	// CODE
	if data["code"] != nil {
		retData["code"] = data["code"]
	} else {
		retData["code"] = "696969"
	}

	// BRANDS
	if data["brands"] != nil {
		retData["brands"] = data["brands"]
	} else {
		retData["brands"] = "MOCK_BRAND"
	}

	// PRODUCT_NAME
	if data["product_name"] != nil {
		retData["product_name"] = data["product_name"]
	} else {
		retData["product_name"] = "MOCK_PROD_NAME"
	}

	// DISPLAY_NAME
	if brands, ok := retData["brands"].(string); ok {
		if prodN, ok2 := retData["product_name"].(string); ok2 {
			retData["display_name"] = brands + prodN
		}
	}

	// categories_hierarchy
	if data["categories_hierarchy"] != nil {
		retData["categories_hierarchy"] = data["categories_hierarchy"]
	} else {
		retData["categories_hierarchy"] = []string{
			"en:mock_main_category",
			"en:mock_sub_category",
		}
	}

	retData["expire_avg"] = "169"
	retData["measurement_unit"] = "liter"
	retData["measurement"] = "69"
	retData["reviewd"] = false

	jsonRetData, err2 := json.Marshal(retData)
	if err2 != nil {
		return nil, err2
	}

	var prod types.Product
	// TODO check if function call is more info
	if unmarshallErr := json.Unmarshal(jsonRetData, &prod); unmarshallErr != nil {
		log.Error(unmarshallErr)
	}

	return &prod, nil
}

func (ai *MockAiCaller) WebScrapeParse(barcode int64) (*types.Product, error) {
	return nil, nil
}

func NewMockAiCaller() *MockAiCaller {
	return &MockAiCaller{}
}
