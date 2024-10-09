package main

import (
	"flag"
	_ "image/png"
	"os"

	"github.com/eiachh/hfc/api"
	"github.com/eiachh/hfc/service"
	"github.com/eiachh/hfc/storage"

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
	aiParser := service.NewAiParser(aiCaller, db)

	hsMan := service.NewHomeStorageManager(db)
	prodMan := service.NewProductManager(db, aiParser)

	prodHandler := api.NewProdHandler(prodMan)
	homestorageHandler := api.NewHsHandler(hsMan, prodMan)

	e := echo.New()
	e.GET("/hs", homestorageHandler.GetAllFood)
	e.POST("/hs/:code", homestorageHandler.AddFood)
	//e.DELETE("/hs/:code", prodHandler.AddFood)

	e.GET("/prod/unverified", prodHandler.GetUnverified)
	e.GET("/prod/categories", prodHandler.GetCatList)
	e.PATCH("/prod/:code", prodHandler.GetUnverified)
	e.GET("/prod/:code", prodHandler.GetProduct)
	e.POST("/prod", prodHandler.NewProd)

	e.Logger.Fatal(e.Start(":1323"))
}

func makeDB() *storage.MongoStorage {
	username := "root"
	password := os.Getenv("MONGODB_ROOT_PASSWORD")
	host := "192.168.49.2"
	port := "30020"
	authDB := "admin"
	offDatabase := "off"
	cacheDatabase := "loc-cache"

	return storage.NewMongoStorage(username, password, host, port, offDatabase, cacheDatabase, authDB)
}
