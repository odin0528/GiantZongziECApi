package frontend

import (
	. "eCommerce/internal/database"
)

type OrderProducts struct {
	ID         int     `json:"-"`
	OrderID    int     `json:"-"`
	ProductID  int     `json:"-"`
	StyleID    int     `json:"-"`
	Qty        int     `json:"qty"`
	Price      float32 `json:"price"`
	Total      float32 `json:"total"`
	Title      string  `json:"title"`
	StyleTitle string  `json:"style_title"`
	Photo      string  `json:"photo"`
	Sku        string  `json:"sku"`
	TimeDefault
}

type OrderProductsCreateReq struct {
	ProductID int                          `json:"productID"`
	Title     string                       `json:"title"`
	Styles    []OrderProductStyleCreateReq `json:"styles"`
}

type OrderProductStyleCreateReq struct {
	StyleID    int     `json:"id"`
	Title      string  `json:"title"`
	StyleTitle string  `json:"style_title"`
	Photo      string  `json:"photo"`
	Sku        string  `json:"sku"`
	Qty        int     `json:"buyCount"`
	Price      float32 `json:"price"`
}

// 基本CURD功能
func (req *OrderProducts) Create() (err error) {
	err = DB.Debug().Create(&req).Error
	return
}
