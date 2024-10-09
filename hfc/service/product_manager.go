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

func (prodMan *ProductManger) GetOrRegisterProduct(barC int64) (*types.Product, error) {
	var prod *types.Product
	var err error

	prod, offJson := prodMan.mongodb.GetByBarCode(barC)

	if prod != nil {
		return prod, nil
	} else if offJson != nil {
		prod, err = prodMan.aiParser.ConvertOffToLocCache(&offJson)
		prod.Code = barC
		prodMan.NewUnReviewed(prod)
	} else {
		prod, err = prodMan.aiParser.DoWebscrape(barC)
		prod.Code = barC
		prodMan.NewUnReviewed(prod)
	}
	return prod, err
}

func (prodMan *ProductManger) Get(barC int64) (*types.Product, error) {
	prod, _ := prodMan.mongodb.GetByBarCode(barC)
	if prod == nil {
		return nil, errors.New("No Product was found with the barcode: " + strconv.FormatInt(barC, 10))
	}
	return prod, nil
}

func (prodMan *ProductManger) GetAllUnreviewed() (*[]types.Product, error) {
	return prodMan.mongodb.GetUnreviewedProducts()
}

func (prodMan *ProductManger) NewReviewed(prod *types.Product) error {
	prod.Reviewed = true
	return prodMan.mongodb.NewProduct(prod)
}

func (prodMan *ProductManger) NewUnReviewed(prod *types.Product) error {
	prod.Reviewed = false
	return prodMan.mongodb.NewProduct(prod)
}

func (prodMan *ProductManger) GetCatListDistinct() (*[]byte, error) {
	return prodMan.mongodb.GetCatListDistinct()
}
