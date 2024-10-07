package api

import (
	"net/http"
	"strconv"

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

// Returns a product based on the given barcode
// TODO needs to check loc-cache and the OFF as well.
func (pHandler *ProdHandler) GetProduct(c echo.Context) error {
	id := c.Param("code")
	barC, convErr := strconv.Atoi(id)
	if convErr != nil {
		return c.String(http.StatusBadRequest, "Barcode has to be int!")
	}

	prod, err := pHandler.prodManager.Get(barC)
	if prod == nil {
		return c.JSON(http.StatusNotFound, err)
	}
	return c.JSON(http.StatusOK, prod)
}

// Adds a new product to the loc-cache db
func (pHandler *ProdHandler) NewProd(c echo.Context) error {

	prod := new(types.Product)
	if err := c.Bind(prod); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := pHandler.prodManager.New(prod); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, prod)
}

func (pHandler *ProdHandler) DeleteProduct(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("code")
	print(id)
	return c.String(http.StatusOK, string("asd"))
}
