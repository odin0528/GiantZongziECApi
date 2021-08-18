package frontend

import (
	. "eCommerce/internal/database"
)

type OrderCreateRequest struct {
	ID           int                      `json:"-"`
	CustomerID   int                      `json:"-"`
	MemberID     int                      `json:"-"`
	Fullname     string                   `json:"fullname"`
	Phone        string                   `json:"phone"`
	Address      string                   `json:"address"`
	Memo         string                   `json:"memo"`
	Method       int                      `json:"method"`
	Total        float32                  `json:"-"`
	Price        float32                  `json:"price"`
	Shipping     float32                  `json:"shipping"`
	Payment      int                      `json:"payment"`
	StoreID      string                   `json:"store_id"`
	StoreName    string                   `json:"store_name"`
	StoreAddress string                   `json:"store_address"`
	StorePhone   string                   `json:"store_phone"`
	Status       int                      `json:"-"`
	Products     []OrderProductsCreateReq `json:"products" gorm:"-"`
	TimeDefault
}

func (OrderCreateRequest) TableName() string {
	return "orders"
}

// 基本CURD功能
func (req *OrderCreateRequest) Create() (err error) {
	req.Total = req.Price + req.Shipping
	err = DB.Debug().Create(&req).Error
	return
}
