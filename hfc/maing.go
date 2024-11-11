package main

import (
	"flag"
	_ "image/png"
	"net/http"
	"os"

	"github.com/eiachh/hfc/api"
	"github.com/eiachh/hfc/service"
	"github.com/eiachh/hfc/storage"

	"github.com/labstack/echo/v4"
)

var (
	// K8s
	healthy bool
	ready   bool

	// Core
	db                 *storage.MongoStorage
	aiCaller           service.AiCaller
	aiParser           *service.AiParser
	hsMan              *service.HomeStorageManager
	prodMan            *service.ProductManger
	prodHandler        *api.ProdHandler
	homestorageHandler *api.HsHandler

	// API
	e    *echo.Echo
	port string
)

// TODO Testing
// TODO MOCKS
// TODO Logging
// TODO Caching vars memory, USE PROD LIST CACHE INSTEAD OF DB, just sync
// TODO Separate logic from mongo storage
func main() {
	ready = false
	healthy = false
	Init()

	e.GET("/hs", homestorageHandler.GetAllFood)
	e.POST("/hs/:code", homestorageHandler.AddFood)
	e.PUT("/hs/:code", homestorageHandler.UpdateFood)
	e.DELETE("/hs/:code", homestorageHandler.DeleteFood)

	e.GET("/prod/unverified", prodHandler.GetUnverified)
	e.GET("/prod/unverified/img/:code", prodHandler.GetUnreviewedImg)
	e.PUT("/prod/unverified/:code", prodHandler.SetUnreviewedImg)
	e.GET("/prod/categories", prodHandler.GetCatList)
	e.GET("/prod/:code", prodHandler.GetProduct)
	e.POST("/prod", prodHandler.NewProd)

	e.GET("/healthy", Healthy)
	e.GET("/ready", Ready)

	ready = true
	healthy = true
	e.Logger.Fatal(e.Start(port))
}

func Healthy(c echo.Context) error {
	if healthy {
		return c.JSON(http.StatusOK, healthy)
	} else {
		return c.JSON(http.StatusServiceUnavailable, healthy)
	}
}

func Ready(c echo.Context) error {
	if ready {
		return c.JSON(http.StatusOK, ready)
	} else {
		return c.JSON(http.StatusServiceUnavailable, ready)
	}
}

func Init() {
	use_openai := flag.String("use_openai", "false", "IF enabled it uses the actual structure like openapi.")
	localDb := flag.String("localdb", "true", "IF enabled it uses a loal db instead of reading everything from env vars. Used when not on helm")
	flag.Parse()

	if *use_openai == "true" {
		aiCaller = service.NewChatGptAiCaller()
	} else {
		aiCaller = service.NewMockAiCaller()
	}

	if *localDb == "true" {
		db = makeDBlocal()
	} else {
		db = makeDBHelm()
	}

	aiParser = service.NewAiParser(aiCaller, db)

	hsMan = service.NewHomeStorageManager(db)
	prodMan = service.NewProductManager(db, aiParser)

	prodHandler = api.NewProdHandler(prodMan)
	homestorageHandler = api.NewHsHandler(hsMan, prodMan)

	port = ":" + os.Getenv("SERVER_PORT")
	e = echo.New()
}

func makeDBlocal() *storage.MongoStorage {
	username := "root"
	password := os.Getenv("MONGODB_PWD")
	host := "192.168.49.2"
	port := "30020"
	authDB := "admin"
	offDatabase := "off"
	cacheDatabase := "loc-cache"

	return storage.NewMongoStorage(username, password, host, port, offDatabase, cacheDatabase, authDB)
}

func makeDBHelm() *storage.MongoStorage {
	username := os.Getenv("MONGODB_USER")
	password := os.Getenv("MONGODB_PWD")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	authDB := os.Getenv("AUTH_DB")
	offDatabase := os.Getenv("OFF_DB")
	cacheDatabase := os.Getenv("CACHE_DB")

	return storage.NewMongoStorage(username, password, host, port, offDatabase, cacheDatabase, authDB)
}
