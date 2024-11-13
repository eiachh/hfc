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
	products            *map[int64]*types.Product
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
	products := make(map[int64]*types.Product)
	allProductAsSlice := *prodMan.mongodb.GetAllProduct()
	for _, prod := range allProductAsSlice {
		products[int64(prod.Code)] = &prod
	}

	prodMan.products = &products
}

func (prodMan *ProductManger) fillCategoryHierarchies() {
	prodMan.categoryHierarchies = prodMan.categoryHierarchies[:0]
	for _, product := range *prodMan.products {
		prodMan.fillCategoryHierarchiesFromProd(product)
	}
}

func (prodMan *ProductManger) fillCategoryHierarchiesFromProd(product *types.Product) {
	var (
		prevCatHItem      *types.CategoryHierarchy
		changedCategories []*types.CategoryHierarchy
	)

	for _, categoryName := range product.Categories {
		cat := prodMan.getCategoryByName(categoryName)
		if cat != nil {
			// Reroute of new hierarchy defined
			if cat.Parent != nil && cat.Parent.Name != prevCatHItem.Name {
				changedCategories = append(changedCategories, cat)
				cat.Parent = prevCatHItem
			}
			prevCatHItem = cat
		} else {
			newItem := types.NewCategoryHierarchyItem(categoryName, prevCatHItem)
			prevCatHItem = newItem
			prodMan.categoryHierarchies = append(prodMan.categoryHierarchies, newItem)
		}
	}

	// Every cat path should be contained once only, change can only happen if the user forcefully changed an existing category list,
	// So we have to update all the products that had the same cat hierarchy
	for _, cat := range changedCategories {
		prodMan.SwitchCategoryInAffectedProducts(*cat)
	}

	if len(changedCategories) > 0 {
		prodMan.fillCategoryHierarchies()
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

func (prodMan *ProductManger) SwitchCategoryInAffectedProducts(changedCategory types.CategoryHierarchy) {
	for _, prod := range *prodMan.products {
		needsChange := false
		for i := 1; i < len(prod.Categories); i++ { // Start from index 1 to avoid out-of-bounds access
			if prod.Categories[i] == changedCategory.Name {
				needsChange = true
				break
			}
		}
		if needsChange {
			copy(prod.Categories[:len(changedCategory.AsSlice())], changedCategory.AsSlice())
			if prod.Reviewed {
				prodMan.NewReviewed(prod)
			} else {
				prodMan.NewUnReviewed(prod)
			}
		}
	}
}

func (prodMan *ProductManger) GetOrRegisterProduct(barC int64) (*types.Product, bool, error) {
	var prod *types.Product
	var err error

	prod, offJson := prodMan.mongodb.GetByBarCode(barC)

	if prod != nil {
		prodMan.fillCategoryHierarchiesFromProd(prod)
		return prod, false, nil
	} else if offJson != nil {
		prod, err = prodMan.aiParser.ConvertOffToLocCache(&offJson)
		if err != nil {
			return nil, false, errors.New("off convert failed")
		}
		prod.Code = barC
		prodMan.NewUnReviewed(prod)
	} else {
		prod, err = prodMan.aiParser.DoWebscrape(barC)
		if err != nil {
			return nil, false, errors.New("webscrape failed")
		}
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
	(*prodMan.products)[prod.Code] = prod
	prodMan.fillCategoryHierarchiesFromProd(prod)

	return prodMan.mongodb.NewProduct(prod)
}

func (prodMan *ProductManger) NewUnReviewed(prod *types.Product) error {
	prod.Reviewed = false
	(*prodMan.products)[prod.Code] = prod
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
