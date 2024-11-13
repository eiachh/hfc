package types

import (
	"time"
)

type StorageItem struct {
	UUID     string    `json:"uuid" bson:"uuid"`
	Acquired time.Time `json:"acquired" bson:"acquired"`
	Expires  time.Time `json:"expires" bson:"expires"`
}

type HomeStorage struct {
	HomeStorageItems map[int64][]StorageItem `json:"home_storage_items" bson:"home_storage_items"`
}

func NewHomeStorage() *HomeStorage {
	return &HomeStorage{
		HomeStorageItems: make(map[int64][]StorageItem),
	}
}
