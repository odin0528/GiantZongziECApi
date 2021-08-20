package frontend

import (
	. "eCommerce/internal/database"
)

type OrderProducts struct {
	ID         int
	OrderID    int
	ProductID  int
	StyleID    int
	Qty        int
	Price      float32
	Total      float32
	Title      string
	StyleTitle string
	Photo      string
	Sku        string
	TimeDefault
}

type OrderProductsCreateReq struct {
	ProductID int                          `json:"productID"`
	Title     string                       `json:"title"`
	Styles    []OrderProductStyleCreateReq `json:"styles"`
}

type OrderProductStyleCreateReq struct {
	StyleID  int     `json:"id"`
	Title    string  `json:"title"`
	SubTitle string  `json:"subTitle"`
	Photo    string  `json:"photo"`
	Sku      string  `json:"sku"`
	Qty      int     `json:"buyCount"`
	Price    float32 `json:"price"`
}

// 基本CURD功能
func (req *OrderProducts) Create() (err error) {
	err = DB.Debug().Create(&req).Error
	return
}
