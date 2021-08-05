package backend

import (
	. "eCommerce/internal/database"
)

type Products struct {
	ID             int              `json:"id"`
	CustomerID     int              `json:"-"`
	Title          string           `json:"title"`
	Description    string           `json:"description"`
	CategoryLayer1 int              `json:"category_layer1"`
	CategoryLayer2 int              `json:"category_layer2"`
	CategoryLayer3 int              `json:"category_layer3"`
	CategoryLayer4 int              `json:"category_layer4"`
	StyleTable     [][]ProductStyle `json:"style_table" gorm:"-"`
	DeletedAt      int              `json:"-"`
	TimeDefault
}

// 基本CURD功能
func (products *Products) Create() (err error) {
	err = DB.Create(&products).Error
	return
}
