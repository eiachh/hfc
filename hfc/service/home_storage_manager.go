package service

import (
	"time"

	"github.com/labstack/gommon/log"

	"github.com/eiachh/hfc/storage"
	"github.com/eiachh/hfc/types"
)

type HomeStorageManager struct {
	homeStorage *types.HomeStorage
	mongodb     *storage.MongoStorage
}

func NewHomeStorageManager(db *storage.MongoStorage) *HomeStorageManager {
	hs, err := getHsFromDb(db)
	if err != nil || hs == nil {
		hs = types.NewHomeStorage()
	}

	if len(hs.Products) < 1 {
		hs = types.NewHomeStorage()
	}

	return &HomeStorageManager{
		homeStorage: hs,
		mongodb:     db,
	}
}

func (manager *HomeStorageManager) InsertProd(amnt int, prod *types.Product) {
	toInsert := types.StorageItem{Accuired: time.Now()}

	for i := 0; i < amnt; i++ {
		manager.homeStorage.Products[prod.Code] = append(manager.homeStorage.Products[prod.Code], toInsert)
	}
	manager.saveHsToDb()
}

func (manager *HomeStorageManager) GetItems(barC int64) []types.StorageItem {
	return manager.homeStorage.Products[barC]
}

func (manager *HomeStorageManager) GetAll() *types.HomeStorage {
	return manager.homeStorage
}

// TODO do not duplicate prod data with save, save the barC only and fetch prod data separately
func (manager *HomeStorageManager) saveHsToDb() {
	err := manager.mongodb.SaveHomeStorage(manager.homeStorage)
	if err != nil {
		log.Error(err)
	}
}

func getHsFromDb(db *storage.MongoStorage) (*types.HomeStorage, error) {
	hs, err := db.LoadHomeStorage()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return hs, nil
}
