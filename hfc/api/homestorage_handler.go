package api

import (
	"net/http"
	"strconv"

	"github.com/eiachh/hfc/service"
	"github.com/eiachh/hfc/types"
	"github.com/labstack/echo/v4"
)

type HsHandler struct {
	HomeStorageManager *service.HomeStorageManager
	ProductMan         *service.ProductManger
}

func NewHsHandler(hs *service.HomeStorageManager, prodMan *service.ProductManger) *HsHandler {
	return &HsHandler{
		HomeStorageManager: hs,
		ProductMan:         prodMan,
	}
}

func (hsh *HsHandler) GetAllFood(c echo.Context) error {
	return c.JSON(http.StatusOK, hsh.HomeStorageManager.GetAll())
}

// Add item to HS endpoint handler
func (hsh *HsHandler) AddFood(c echo.Context) error {
	barC, convErr := strconv.ParseInt(c.Param("code"), 10, 64)
	if convErr != nil {
		return c.JSON(http.StatusBadRequest, convErr)
	}
	var requestBody struct {
		Amount int `json:"amount"`
	}
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request body",
		})
	}

	reqImg, err := hsh.hsAddFood(barC, requestBody.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if reqImg {
		return c.JSON(http.StatusPartialContent, "require img for verification")
	}

	return c.JSON(http.StatusOK, barC)
}

func (hsh *HsHandler) UpdateFood(c echo.Context) error {
	barC, convErr := strconv.ParseInt(c.Param("code"), 10, 64)
	if convErr != nil {
		return c.JSON(http.StatusBadRequest, convErr)
	}
	var requestBodyAsStorItem *types.StorageItem
	if err := c.Bind(&requestBodyAsStorItem); err != nil || requestBodyAsStorItem == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request body",
		})
	}

	err := hsh.HomeStorageManager.UpdateItem(barC, requestBodyAsStorItem)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, "ok")
}

func (hsh *HsHandler) DeleteFood(c echo.Context) error {
	barC, convErr := strconv.ParseInt(c.Param("code"), 10, 64)
	if convErr != nil {
		return c.JSON(http.StatusBadRequest, convErr)
	}
	var requestBody struct {
		Uuid string `json:"uuid"`
	}
	if err := c.Bind(&requestBody); err != nil || requestBody.Uuid == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request body",
		})
	}

	uuid := requestBody.Uuid
	err := hsh.HomeStorageManager.RemoveItem(barC, uuid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, "ok")
}

func (hsh *HsHandler) hsAddFood(barC int64, amnt int) (bool, error) {
	//MOCK
	return true, nil
	//MOCK

	prod, reqImg, err := hsh.ProductMan.GetOrRegisterProduct(barC)
	if err != nil {
		return false, err
	}

	hsh.HomeStorageManager.InsertProd(amnt, prod)
	return reqImg, nil
}
