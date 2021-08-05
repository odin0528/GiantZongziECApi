package backend

import (
	. "eCommerce/internal/database"
)

type ProductStyle struct {
	ID         int    `json:"id"`
	CustomerID int    `json:"-"`
	ProductID  int    `json:"-"`
	Title      string `json:"title"`
	SubTitle   string `json:"subTitle"`
	Sku        string `json:"sku"`
	Price      int    `json:"price"`
	Qty        int    `json:"qty"`
	DeletedAt  int    `json:"-"`
	TimeDefault
}

func (ProductStyle) TableName() string {
	return "product_style"
}

// 基本CURD功能
func (style *ProductStyle) Create() (err error) {
	err = DB.Create(&style).Error
	return
}
