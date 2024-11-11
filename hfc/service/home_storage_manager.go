package service

import (
	"errors"
	"time"

	"github.com/google/uuid"
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

	if len(hs.HomeStorageItems) < 1 {
		hs = types.NewHomeStorage()
	}

	return &HomeStorageManager{
		homeStorage: hs,
		mongodb:     db,
	}
}

func (manager *HomeStorageManager) InsertProd(amnt int, prod *types.Product) {
	toInsert := types.StorageItem{
		Acquired: time.Now(),
		Expires:  time.Now().Add(time.Duration(prod.ExpireDays) * 24 * time.Hour),
	}

	for i := 0; i < amnt; i++ {
		toInsert.UUID = uuid.NewString()
		manager.homeStorage.HomeStorageItems[prod.Code] = append(manager.homeStorage.HomeStorageItems[prod.Code], toInsert)
	}
	manager.saveHsToDb()
}

func (manager *HomeStorageManager) UpdateItem(barCode int64, receivedItem *types.StorageItem) error {
	itemlist, hasKey := manager.homeStorage.HomeStorageItems[barCode]
	if !hasKey {
		return errors.New("item was not present to remove")
	}

	for ind, item := range itemlist {
		if item.UUID == receivedItem.UUID {
			itemlist[ind] = *receivedItem
			manager.saveHsToDb()
			return nil
		}
	}
	return errors.New("item was not present to remove")
}

func (manager *HomeStorageManager) RemoveItem(barCode int64, uuid string) error {
	itemlist, hasKey := manager.homeStorage.HomeStorageItems[barCode]
	if !hasKey {
		return errors.New("item was not present to remove")
	}
	indToRemove := -1
	for ind, item := range itemlist {
		if item.UUID == uuid {
			indToRemove = ind
			break
		}
	}
	if indToRemove == -1 {
		return errors.New("item was not present to remove")
	}
	if len(itemlist) == 1 {
		delete(manager.homeStorage.HomeStorageItems, barCode)
	} else {
		itemlist[indToRemove] = itemlist[len(itemlist)-1]
		itemlist = itemlist[:len(itemlist)-1]
		manager.homeStorage.HomeStorageItems[barCode] = itemlist
	}
	manager.saveHsToDb()
	return nil
}

func (manager *HomeStorageManager) GetItems(barC int64) []types.StorageItem {
	return manager.homeStorage.HomeStorageItems[barC]
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
