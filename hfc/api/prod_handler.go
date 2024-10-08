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

// Returns a product based on the given barcode from loc-cache
func (pHandler *ProdHandler) GetProduct(c echo.Context) error {
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
	prod, err := pHandler.prodManager.GetAllUnreviewed()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if len(*prod) == 0 {
		return c.JSON(http.StatusNotFound, "No unreviewed item was found!")
	}
	return c.JSON(http.StatusOK, prod)
}

// Adds a new product to the loc-cache db
func (pHandler *ProdHandler) NewProd(c echo.Context) error {

	prod := new(types.Product)
	if err := c.Bind(prod); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := pHandler.prodManager.NewReviewed(prod); err != nil {
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

func (pHandler *ProdHandler) GetCatList(c echo.Context) error {
	jsonCat, err := pHandler.prodManager.GetCatListDistinct()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, string(*jsonCat))
}
