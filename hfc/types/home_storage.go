package types

import (
	"time"
)

type StorageItem struct {
	Accuired time.Time
}

// TODO separate by user
type HomeStorage struct {
	Products map[int64][]StorageItem
}

func NewHomeStorage() *HomeStorage {
	return &HomeStorage{
		Products: make(map[int64][]StorageItem),
	}
}
