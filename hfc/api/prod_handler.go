package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/eiachh/hfc/service"
	"github.com/eiachh/hfc/storage"
	"github.com/eiachh/hfc/types"
	"github.com/labstack/echo/v4"
)

type ProdHandler struct {
	mongodb     *storage.MongoStorage
	homestorage *types.HomeStorage
	aiParser    *service.AiParser
}

func NewProdHandler(db *storage.MongoStorage, hs *types.HomeStorage, ai *service.AiParser) *ProdHandler {
	return &ProdHandler{
		mongodb:     db,
		homestorage: hs,
		aiParser:    ai,
	}
}

// Returns a product based on the given barcode
// TODO needs to check loc-cache and the OFF as well.
func (pHandler *ProdHandler) GetFood(c echo.Context) error {
	id := c.Param("code")
	prod, _ := pHandler.mongodb.GetByBarCode(id)
	if prod == nil {
		return c.String(http.StatusNotFound, "No product was found in processed cache")
	}
	prodJson, _ := json.Marshal(prod)
	return c.String(http.StatusOK, string(prodJson))
}

// Add item to HS endpoint handler
func (pHandler *ProdHandler) AddFood(c echo.Context) error {
	barC, _ := strconv.Atoi(c.Param("code"))
	var requestBody struct {
		Amount string `json:"amount"`
	}

	// Bind the request body to the struct
	if err := c.Bind(&requestBody); err != nil || requestBody.Amount == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request body",
		})
	}
	amount := requestBody.Amount

	// TODO amnt has to be str bruh
	err := pHandler.handleHSAdd(barC, amount, c)
	if err != nil {
		return nil
	}

	return c.String(http.StatusOK, "Added "+strconv.Itoa(barC)+", amnt: "+amount)
}

// Adds a new product to the loc-cache db
func (pHandler *ProdHandler) NewFood(c echo.Context) error {

	prod := new(types.Product)
	if err := c.Bind(prod); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if res := pHandler.mongodb.New(prod); res {
		return c.JSON(http.StatusOK, prod)
	}
	return c.JSON(http.StatusInternalServerError, "Could not insert to DB")
}

// Recives a barcode from the scanner or the image TBD its MOCK
func (pHandler *ProdHandler) MOCKScanBarcode(c echo.Context) error {
	amnt := c.Param("amnt")
	return c.String(http.StatusOK, amnt)
}

// // Recieves a img of a barcode, converts it to the int barcode and adds it to the home-storage
// func (pHandler *ProdHandler) PhoneBarcode(c echo.Context) error {
// 	amnt := c.Param("amnt")

// 	//barC := parser.MOCKConvertBarcodeImg(nil)
// 	return pHandler.handleHSAdd(barC, amnt, c)
// }

func (pHandler *ProdHandler) handleHSAdd(barC int, amnt string, c echo.Context) error {
	prod, offJson := pHandler.mongodb.GetByBarCode(strconv.Itoa(barC))

	if prod != nil {
		pHandler.insertProdToHs(amnt, prod)
		return c.JSON(http.StatusNotFound, "Missing barcode from DB")
	} else if offJson != nil {
		pHandler.aiParser.ConvertOffToLocCache(&offJson)
	} else {

	}

	return nil
}

func (pHandler *ProdHandler) insertProdToHs(amnt string, prod *types.Product) {
	amntInt, err := strconv.Atoi(amnt)
	if err != nil || amnt == "" {
		pHandler.homestorage.InsertProd(1, prod)
	} else {
		pHandler.homestorage.InsertProd(amntInt, prod)
	}
}

func (pHandler *ProdHandler) DeleteFood(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("code")
	print(id)
	return c.String(http.StatusOK, string("asd"))
}
