package main

import (
	"flag"
	_ "image/png"
	"net/http"

	"github.com/eiachh/hfc/api"
	"github.com/eiachh/hfc/service"
	"github.com/eiachh/hfc/storage"
	"github.com/eiachh/hfc/types"

	"github.com/labstack/echo/v4"
)

func main() {
	prod_enabled := flag.String("prod_enabled", "false", "IF enabled it uses the actual structure like openapi.")
	flag.Parse()

	var aiCaller service.AiCaller
	if *prod_enabled == "true" {
		aiCaller = service.NewChatGptAiCaller()
	} else {
		aiCaller = service.NewMockAiCaller()
	}

	db := makeDB()
	prodHandler := api.NewProdHandler(db, types.NewHomeStorage(), service.NewAiParser(aiCaller))
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/prod/get/:code", prodHandler.GetFood)
	e.POST("/prod/post/:code", prodHandler.AddFood)
	e.POST("/prod/new", prodHandler.NewFood)
	e.POST("/prod/scan/code/:amnt", prodHandler.MOCKScanBarcode)

	e.Logger.Fatal(e.Start(":1323"))
}

func makeDB() *storage.MongoStorage {
	username := "root"
	// TODO Un-dox yourself 4head
	password := "lDyd8IubHC"
	host := "192.168.49.2"
	port := "30020"
	authDB := "admin"
	offDatabase := "off"
	cacheDatabase := "loc-cache"

	return storage.NewMongoStorage(username, password, host, port, offDatabase, cacheDatabase, authDB)
}
