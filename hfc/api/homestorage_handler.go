package api

import (
	"net/http"
	"strconv"

	"github.com/eiachh/hfc/service"
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
		Amount string `json:"amount"`
	}
	if err := c.Bind(&requestBody); err != nil || requestBody.Amount == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "Invalid request body",
		})
	}

	amount := requestBody.Amount
	amntInt, convErr := strconv.Atoi(amount)
	if convErr != nil {
		amntInt = 1
	}
	err := hsh.hsAddFood(barC, amntInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.String(http.StatusOK, "Added "+strconv.FormatInt(barC, 10)+", amnt: "+amount)
}

func (hsh *HsHandler) hsAddFood(barC int64, amnt int) error {
	prod, err := hsh.ProductMan.GetOrRegisterProduct(barC)
	if err != nil {
		return err
	}

	hsh.HomeStorageManager.InsertProd(amnt, prod)
	return nil
}
