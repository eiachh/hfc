package service

import "github.com/eiachh/hfc/types"

type AiCaller interface {
	ParseOff([]byte) (*types.Product, error)
	WebScrapeParse(barcode int64) (*types.Product, error)
}
