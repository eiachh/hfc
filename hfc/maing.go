package main

import (
	_ "image/png"
	"net/http"

	"github.com/eiachh/hfc/api"
	"github.com/eiachh/hfc/storage"
	"github.com/eiachh/hfc/types"

	"github.com/labstack/echo/v4"
)

func main() {

	db := makeDB()
	prodHandler := api.NewProdHandler(db, types.NewHomeStorage())
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/prod/get/:code", prodHandler.GetFood)
	e.POST("/prod/post/:code", prodHandler.AddFood)
	e.POST("/prod/new", prodHandler.NewFood)
	e.POST("/prod/phone/code/:amnt", prodHandler.PhoneBarcode)
	e.POST("/prod/scan/code/:amnt", prodHandler.MOCKScanBarcode)

	e.Logger.Fatal(e.Start(":1323"))
}

func makeDB() *storage.MongoStorage {
	username := "root"
	password := "lDyd8IubHC"
	host := "192.168.49.2"
	port := "30020"
	authDB := "admin"       // Authentication database
	database := "loc-cache" // Replace with your database name

	return storage.NewMongoStorage(username, password, host, port, database, authDB)
}
