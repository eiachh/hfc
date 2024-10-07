package types

import (
	"encoding/json"
	"strconv"
)

type ProductStr struct {
	Code            string   `json:"code" bson:"code"`
	Brands          string   `json:"brands" bson:"brands"`
	Name            string   `json:"product_name" bson:"product_name"`
	DisplayName     string   `json:"display_name" bson:"display_name"`
	Expire          string   `json:"expire_avg" bson:"expire_avg"`
	Categories      []string `json:"categories_hierarchy" bson:"categories_hierarchy"`
	Reviewed        string   `json:"reviewed" bson:"reviewed"`
	MeasurementUnit string   `json:"measurement_unit" bson:"measurement_unit"`
}

type Product struct {
	Code            int      `json:"code" bson:"code"`
	Brands          string   `json:"brands" bson:"brands"`
	Name            string   `json:"product_name" bson:"product_name"`
	DisplayName     string   `json:"display_name" bson:"display_name"`
	Expire          int      `json:"expire_avg" bson:"expire_avg"`
	Categories      []string `json:"categories_hierarchy" bson:"categories_hierarchy"`
	Reviewed        bool     `json:"reviewed" bson:"reviewed"`
	MeasurementUnit string   `json:"measurement_unit" bson:"measurement_unit"`
}

func ProdWithCode(code int) *Product {
	return &Product{
		Code:     code,
		Reviewed: false,
	}
}

func (p *Product) AsStr() *ProductStr {
	return &ProductStr{
		// Convert int Code to string
		Code:        strconv.Itoa(p.Code),
		Brands:      p.Brands,
		Name:        p.Name,
		DisplayName: p.DisplayName,
		// Format time.Time Expire to string (YYYY-MM-DD)
		Expire:     strconv.Itoa(p.Expire),
		Categories: p.Categories,
		Reviewed:   strconv.FormatBool(p.Reviewed),
	}

}

// Custom Unmarshal for Product to handle code as string
func (p *Product) UnmarshalJSON(data []byte) error {
	var aux struct {
		Code        string   `json:"code"` // Temporary string for unmarshalling
		Brands      string   `json:"brands"`
		Name        string   `json:"product_name"`
		DisplayName string   `json:"display_name"`
		Expire      string   `json:"expire_avg"` // Temporarily as string
		Categories  []string `json:"categories_hierarchy"`
		Reviewed    string   `json:"reviewed" bson:"reviewed"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Convert code string to int
	if aux.Code != "" {
		code, err := strconv.Atoi(aux.Code)
		if err != nil {
			return err
		}
		p.Code = code
	}

	p.Brands = aux.Brands
	p.Name = aux.Name
	p.DisplayName = aux.DisplayName

	// Parse expire string to time.Time
	expire, err := strconv.Atoi(aux.Expire)
	if err != nil {
		return err
	}
	p.Expire = expire
	p.Categories = aux.Categories
	p.Reviewed = aux.Reviewed == "true"

	return nil
}
