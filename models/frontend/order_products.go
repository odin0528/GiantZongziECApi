package frontend

import (
	. "eCommerce/internal/database"
)

type OrderProducts struct {
	ID              int     `json:"-"`
	OrderID         int     `json:"-"`
	ProductID       int     `json:"-"`
	StyleID         int     `json:"-"`
	Qty             int     `json:"qty"`
	Price           float64 `json:"price"`
	IsDiscount      bool    `json:"is_discount"`
	Discount        float64 `json:"discount"`
	DiscountedPrice float64 `json:"discounted_price"`
	Total           float64 `json:"total"`
	Title           string  `json:"title"`
	StyleTitle      string  `json:"style_title"`
	Photo           string  `json:"photo"`
	Sku             string  `json:"sku"`
	TimeDefault
}

type OrderProductsCreateReq struct {
	ProductID int                          `json:"productID"`
	Title     string                       `json:"title"`
	Styles    []OrderProductStyleCreateReq `json:"styles"`
}

type OrderProductStyleCreateReq struct {
	StyleID         int     `json:"id"`
	Title           string  `json:"title"`
	StyleTitle      string  `json:"style_title"`
	Photo           string  `json:"photo"`
	Sku             string  `json:"sku"`
	BuyCount        int     `json:"buy_count"`
	Qty             int     `json:"qty"`
	Price           float64 `json:"price"`
	IsDiscount      bool    `json:"is_discount"`
	Discount        float64 `json:"discount"`
	DiscountedPrice float64 `json:"discounted_price"`
	NoStoreDelivery int     `json:"no_store_delivery"`
	NoOverSale      bool    `json:"no_over_sale"`
}

// 基本CURD功能
func (req *OrderProducts) Create() (err error) {
	err = DB.Create(&req).Error
	return
}
