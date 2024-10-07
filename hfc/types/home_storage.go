package types

import (
	"time"
)

type StorageItem struct {
	Accuired time.Time
	Product  *Product
}

type HomeStorage struct {
	Products map[int][]StorageItem
}

func NewHomeStorage() *HomeStorage {
	return &HomeStorage{
		Products: make(map[int][]StorageItem),
	}
}
