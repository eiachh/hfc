package service

import (
	"errors"
	"strconv"

	"github.com/eiachh/hfc/storage"
	"github.com/eiachh/hfc/types"
)

type ProductManger struct {
	mongodb             *storage.MongoStorage
	aiParser            *AiParser
	products            []types.Product
	categoryHierarchies []*types.CategoryHierarchy
}

func NewProductManager(db *storage.MongoStorage, ai *AiParser) *ProductManger {
	prodman := &ProductManger{
		mongodb:  db,
		aiParser: ai,
	}
	prodman.init()
	return prodman
}

func (prodMan *ProductManger) init() {
	prodMan.fillProducts()
	prodMan.fillCategoryHierarchies()
}

func (prodMan *ProductManger) fillProducts() {
	prodMan.products = *prodMan.mongodb.GetAllProduct()
}

func (prodMan *ProductManger) fillCategoryHierarchies() {
	for _, product := range prodMan.products {
		prodMan.fillCategoryHierarchiesFromProd(&product)
	}
}

func (prodMan *ProductManger) fillCategoryHierarchiesFromProd(product *types.Product) {
	var prevCatHItem *types.CategoryHierarchy
	for _, categoryName := range product.Categories {
		cat := prodMan.getCategoryByName(categoryName)
		if cat != nil {
			// Reroute of new hierarchy defined
			if cat.Parent != nil && cat.Parent.Name != prevCatHItem.Name {
				cat.Parent = prevCatHItem
			}
			prevCatHItem = cat
		} else {
			newItem := types.NewCategoryHierarchyItem(categoryName, prevCatHItem)
			prevCatHItem = newItem
			prodMan.categoryHierarchies = append(prodMan.categoryHierarchies, newItem)
		}
	}
}

func (prodMan *ProductManger) getCategoryByName(nameOfHierarchyItem string) *types.CategoryHierarchy {
	for _, catHItem := range prodMan.categoryHierarchies {
		if catHItem.Name == nameOfHierarchyItem {
			return catHItem
		}
	}
	return nil
}

// TODO cat hierarchy
func (prodMan *ProductManger) GetOrRegisterProduct(barC int64) (*types.Product, bool, error) {
	var prod *types.Product
	var err error

	prod, offJson := prodMan.mongodb.GetByBarCode(barC)

	if prod != nil {
		prodMan.fillCategoryHierarchiesFromProd(prod)
		return prod, false, nil
	} else if offJson != nil {
		prod, err = prodMan.aiParser.ConvertOffToLocCache(&offJson)
		prod.Code = barC
		prodMan.NewUnReviewed(prod)
	} else {
		prod, err = prodMan.aiParser.DoWebscrape(barC)
		prod.Code = barC
		prodMan.NewUnReviewed(prod)
	}
	prodMan.fillCategoryHierarchiesFromProd(prod)
	return prod, true, err
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
	prodMan.fillCategoryHierarchiesFromProd(prod)
	return prodMan.mongodb.NewProduct(prod)
}

func (prodMan *ProductManger) NewUnReviewed(prod *types.Product) error {
	prod.Reviewed = false
	prodMan.fillCategoryHierarchiesFromProd(prod)
	return prodMan.mongodb.NewProduct(prod)
}

func (prodMan *ProductManger) SetUnreviewedImg(barCode int64, imgAsBase64 string) error {
	return prodMan.mongodb.NewUnreviewedImg(barCode, imgAsBase64)
}

func (prodMan *ProductManger) GetUnreviewedImg(barCode int64) (string, error) {
	return prodMan.mongodb.GetUnreviewedImg(barCode)
}

func (prodMan *ProductManger) GetCatListDistinct() ([]string, error) {
	var retList []string
	for _, categoryHItem := range prodMan.categoryHierarchies {
		retList = append(retList, categoryHItem.ToString())
	}
	return retList, nil
}
