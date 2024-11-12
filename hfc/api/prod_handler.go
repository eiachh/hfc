package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/eiachh/hfc/logger"
	"github.com/eiachh/hfc/service"
	"github.com/eiachh/hfc/types"

	"github.com/labstack/echo/v4"
)

type ProdHandler struct {
	prodManager *service.ProductManger
}

func NewProdHandler(pMan *service.ProductManger) *ProdHandler {
	return &ProdHandler{
		prodManager: pMan,
	}
}

// Returns a product based on the given barcode from loc-cache
func (pHandler *ProdHandler) GetProduct(c echo.Context) error {
	logger.Log().Info("GetProduct called")
	id := c.Param("code")
	barC, convErr := strconv.ParseInt(id, 10, 64)
	if convErr != nil {
		return c.String(http.StatusBadRequest, "Barcode has to be int!")
	}

	prod, err := pHandler.prodManager.Get(barC)
	if prod == nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, prod)
}

func (pHandler *ProdHandler) GetUnverified(c echo.Context) error {
	logger.Log().Info("GetUnverified called")
	prods, err := pHandler.prodManager.GetAllUnreviewed()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if len(*prods) == 0 {
		return c.JSON(http.StatusNotFound, "No unreviewed item was found!")
	}
	return c.JSON(http.StatusOK, prods)
}

func (pHandler *ProdHandler) SetUnreviewedImg(c echo.Context) error {
	logger.Log().Info("GetUnverified called")
	barC, convErr := strconv.ParseInt(c.Param("code"), 10, 64)
	if convErr != nil {
		return c.JSON(http.StatusBadRequest, convErr)
	}
	var requestBody struct {
		ImgAsBase64 string `json:"imgAsBase64"`
	}
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request body",
		})
	}

	reqBodyJson, _ := json.Marshal(requestBody)
	logger.Log().Debugf("Request with barC: %d, requestBody: %s", barC, reqBodyJson)
	if err := pHandler.prodManager.SetUnreviewedImg(barC, requestBody.ImgAsBase64); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "OK")
}

func (pHandler *ProdHandler) GetUnreviewedImg(c echo.Context) error {
	logger.Log().Info("GetUnreviewedImg called")
	barC, convErr := strconv.ParseInt(c.Param("code"), 10, 64)
	if convErr != nil {
		return c.JSON(http.StatusBadRequest, convErr)
	}

	img64, err := pHandler.prodManager.GetUnreviewedImg(barC)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, img64)
}

// Adds a new product to the loc-cache db
func (pHandler *ProdHandler) NewProd(c echo.Context) error {
	logger.Log().Info("NewProd called")
	prod := new(types.Product)
	if err := c.Bind(prod); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	reqBodyJson, _ := json.Marshal(prod)
	logger.Log().Debugf("Request with requestBody: %s", reqBodyJson)
	if err := pHandler.prodManager.NewReviewed(prod); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, prod)
}

func (pHandler *ProdHandler) GetCatList(c echo.Context) error {
	logger.Log().Info("GetCatList called")
	catList, err := pHandler.prodManager.GetCatListDistinct()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, catList)
}
