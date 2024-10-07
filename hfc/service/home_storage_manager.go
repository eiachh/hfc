package service

import (
	"time"

	"github.com/eiachh/hfc/storage"
	"github.com/eiachh/hfc/types"
)

type HomeStorageManager struct {
	homeStorage *types.HomeStorage
	mongodb     *storage.MongoStorage
}

func NewHomeStorageManager(db *storage.MongoStorage) *HomeStorageManager {
	return &HomeStorageManager{
		homeStorage: types.NewHomeStorage(),
		mongodb:     db,
	}
}

func (manager *HomeStorageManager) InsertProd(amnt int, prod *types.Product) {
	toInsert := types.StorageItem{Product: prod, Accuired: time.Now()}

	for i := 0; i < amnt; i++ {
		manager.homeStorage.Products[prod.Code] = append(manager.homeStorage.Products[prod.Code], toInsert)
	}
}

func (s *HomeStorageManager) GetItems(barC int) []types.StorageItem {
	return s.homeStorage.Products[barC]
}

func (s *HomeStorageManager) GetAll() *types.HomeStorage {
	return s.homeStorage
}

func (s *HomeStorageManager) SaveHsToDb() error {
	return nil
}
