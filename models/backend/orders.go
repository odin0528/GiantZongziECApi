package backend

import (
	. "eCommerce/internal/database"
)

type OrderListReq struct {
	CustomerID int `json:"-"`
	Pagination
}

type Orders struct {
	ID           int             `json:"id"`
	CustomerID   int             `json:"-"`
	MemberID     int             `json:"-"`
	Fullname     string          `json:"fullname"`
	Phone        string          `json:"phone"`
	Address      string          `json:"address"`
	Memo         string          `json:"memo"`
	Method       int             `json:"method"`
	Total        float32         `json:"total"`
	Price        float32         `json:"price"`
	Shipping     float32         `json:"shipping"`
	Payment      int             `json:"payment"`
	StoreID      string          `json:"store_id"`
	StoreName    string          `json:"store_name"`
	StoreAddress string          `json:"store_address"`
	StorePhone   string          `json:"store_phone"`
	Status       int             `json:"status"`
	Products     []OrderProducts `json:"products" gorm:"-"`
	TimeDefault
}

// 查詢功能
func (req *OrderListReq) FetchAll() (orders []Orders, pagination Pagination) {
	var count int64
	sql := DB.Debug().Model(&Orders{}).Where("customer_id = ?", req.CustomerID)
	sql.Count(&count)
	sql.Offset((req.Page - 1) * req.Items).Limit(req.Items).Scan(&orders)
	pagination = CreatePagination(req.Page, req.Items, count)
	return
}

// 關連功能
func (order *Orders) GetProducts() {
	DB.Model(&OrderProducts{}).Where("order_id = ?", order.ID).Order("created_at ASC").Scan(&order.Products)
}
