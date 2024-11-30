package service

import "github.com/eiachh/hfc/types"

type AiCaller interface {
	ParseOff([]byte) (*types.Product, error)
	WebScrapeParse(barcode int64, chatComp *types.ChatCompletion, aibody *types.AiReqBody, callCount int) (*types.Product, error)
}
