package types

import (
	"time"
)

type StorageItem struct {
	Accuired time.Time
}

// TODO separate by user
type HomeStorage struct {
	HomeStorageItems map[int64][]StorageItem
}

func NewHomeStorage() *HomeStorage {
	return &HomeStorage{
		HomeStorageItems: make(map[int64][]StorageItem),
	}
}
