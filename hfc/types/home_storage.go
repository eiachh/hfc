package types

type StorageItem struct {
	Ammount int
	Product *Product
}

type HomeStorage struct {
	Products map[int]StorageItem
}

func NewHomeStorage() *HomeStorage {
	return &HomeStorage{
		Products: make(map[int]StorageItem),
	}
}

func (s *HomeStorage) InsertProd(amnt int, prod *Product) {
	if item, exists := s.Products[prod.Code]; exists {
		item.Ammount += amnt
	} else {
		s.Products[prod.Code] = StorageItem{Ammount: amnt, Product: prod}
	}
}
