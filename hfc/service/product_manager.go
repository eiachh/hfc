package service

import (
	"errors"
	"strconv"

	"github.com/eiachh/hfc/storage"
	"github.com/eiachh/hfc/types"
)

type ProductManger struct {
	mongodb  *storage.MongoStorage
	aiParser *AiParser
}

func NewProductManager(db *storage.MongoStorage, ai *AiParser) *ProductManger {
	return &ProductManger{
		mongodb:  db,
		aiParser: ai,
	}
}

func (prodMan *ProductManger) GetOrRegisterProduct(barC int) (*types.Product, error) {
	var prod *types.Product
	var err error

	prod, offJson := prodMan.mongodb.GetByBarCode(barC)

	if prod != nil {
		return prod, nil
	} else if offJson != nil {
		prod, err = prodMan.aiParser.ConvertOffToLocCache(&offJson)
	} else {
		prod, err = prodMan.aiParser.DoWebscrape(barC)
	}
	return prod, err
}

func (prodMan *ProductManger) Get(barC int) (*types.Product, error) {
	prod, _ := prodMan.mongodb.GetByBarCode(barC)
	if prod == nil {
		return nil, errors.New("No Product was found with the barcode: " + strconv.Itoa(barC))
	}
	return prod, nil
}

func (prodMan *ProductManger) New(prod *types.Product) error {
	return prodMan.mongodb.NewProduct(prod)
}
