package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/eiachh/hfc/logger"
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
	logger.Log().Info("GetAllFood called")
	return c.JSON(http.StatusOK, hsh.HomeStorageManager.GetAll())
}

// Add item to HS endpoint handler
func (hsh *HsHandler) AddFood(c echo.Context) error {
	logger.Log().Info("AddFood called")
	barC, convErr := strconv.ParseInt(c.Param("code"), 10, 64)
	if convErr != nil {
		return c.JSON(http.StatusBadRequest, convErr)
	}
	var requestBody struct {
		Amount int `json:"amount"`
	}
	if err := c.Bind(&requestBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("Invalid request body. %s", err),
		})
	}

	reqBodyJson, _ := json.Marshal(requestBody)
	logger.Log().Debugf("Request with barC: %d, requestBody: %s", barC, reqBodyJson)
	reqImg, err := hsh.hsAddFood(barC, requestBody.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	if reqImg {
		return c.JSON(http.StatusPartialContent, map[string]string{
			"message": "require img for verification",
		})
	}

	return c.JSON(http.StatusOK, map[string]int64{
		"message": barC,
	})
}

func (hsh *HsHandler) UpdateFood(c echo.Context) error {
	logger.Log().Info("UpdateFood called")
	barC, convErr := strconv.ParseInt(c.Param("code"), 10, 64)
	if convErr != nil {
		return c.JSON(http.StatusBadRequest, convErr)
	}

	var requestBodyAsStorItem *types.StorageItem
	if err := c.Bind(&requestBodyAsStorItem); err != nil || requestBodyAsStorItem == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("Invalid request body. %s", err),
		})
	}

	reqBodyJson, _ := json.Marshal(requestBodyAsStorItem)
	logger.Log().Debugf("Request with barC: %d, requestBody: %s", barC, reqBodyJson)
	err := hsh.HomeStorageManager.UpdateItem(barC, requestBodyAsStorItem)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "ok",
	})
}

func (hsh *HsHandler) DeleteFood(c echo.Context) error {
	logger.Log().Info("DeleteFood called")
	barC, convErr := strconv.ParseInt(c.Param("code"), 10, 64)
	if convErr != nil {
		return c.JSON(http.StatusBadRequest, convErr)
	}
	var requestBody struct {
		Uuid string `json:"uuid"`
	}
	if err := c.Bind(&requestBody); err != nil || requestBody.Uuid == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("Invalid request body. %s", err),
		})
	}

	reqBodyJson, _ := json.Marshal(requestBody)
	logger.Log().Debugf("Request with barC: %d, requestBody: %s", barC, reqBodyJson)
	uuid := requestBody.Uuid
	err := hsh.HomeStorageManager.RemoveItem(barC, uuid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, map[string]string{
		"message": "ok",
	})
}

func (hsh *HsHandler) hsAddFood(barC int64, amnt int) (bool, error) {
	prod, reqImg, err := hsh.ProductMan.GetOrRegisterProduct(barC)
	if err != nil {
		return false, err
	}

	hsh.HomeStorageManager.InsertProd(amnt, prod)
	return reqImg, nil
}
