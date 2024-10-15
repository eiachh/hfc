package types

type Product struct {
	Code            int64    `json:"code" bson:"code"`
	Brands          string   `json:"brands" bson:"brands"`
	Name            string   `json:"product_name" bson:"product_name"`
	DisplayName     string   `json:"display_name" bson:"display_name"`
	ExpireDays      int      `json:"expire_avg" bson:"expire_avg"`
	Categories      []string `json:"categories_hierarchy" bson:"categories_hierarchy"`
	Reviewed        bool     `json:"reviewed" bson:"reviewed"`
	MeasurementUnit string   `json:"measurement_unit" bson:"measurement_unit"`
}

type ALTProduct struct {
	Code            int      `json:"code" bson:"code"`
	Brands          string   `json:"brands" bson:"brands"`
	Name            string   `json:"product_name" bson:"product_name"`
	DisplayName     string   `json:"display_name" bson:"display_name"`
	Expire          int      `json:"expire_avg" bson:"expire_avg"`
	Categories      []string `json:"categories_hierarchy" bson:"categories_hierarchy"`
	Reviewed        bool     `json:"reviewed" bson:"reviewed"`
	MeasurementUnit string   `json:"measurement_unit" bson:"measurement_unit"`
}
