package service

import (
	"errors"
	"strings"

	"github.com/eiachh/hfc/logger"
	"github.com/eiachh/hfc/types"
)

type MockAiCaller struct{}

func (mockAi *MockAiCaller) ParseOff(trimmedOff []byte) (*types.Product, error) {
	logger.Log().Warn("You used the mock ai FOREHEAD")
	if strings.Contains(string(trimmedOff), "4014500513010") {
		prod := &types.Product{
			Brands:          "Zott",
			DisplayName:     "Jogobella Forest Fruit Yogurt",
			ExpireDays:      10,
			Categories:      []string{"dairies", "fermented-foods", "fermented-milk-products", "desserts", "dairy-desserts", "fermented-dairy-desserts", "fermented-dairy-desserts-with-fruits", "yogurts", "fruit-yogurts"},
			Reviewed:        false,
			MeasurementUnit: "gramm",
		}
		return prod, nil
	}
	return nil, errors.New("did not contain yogobella barcode")
}

func (ai *MockAiCaller) WebScrapeParse(barcode int64) (*types.Product, error) {
	logger.Log().Warn("You used the mock ai FOREHEAD")
	return nil, errors.New("ou used the mock ai FOREHEAD")
}

func NewMockAiCaller() *MockAiCaller {
	logger.Log().Warn("You used the mock ai FOREHEAD")
	return &MockAiCaller{}
}
