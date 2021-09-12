package frontend

import (
	. "eCommerce/internal/database"

	"gorm.io/gorm"
)

type OrderQuery struct {
	MemberID int `json:"-"`
	Pagination
}

type OrderCreateRequest struct {
	ID           int                      `json:"-"`
	PlatformID   int                      `json:"-"`
	MemberID     int                      `json:"-"`
	Email        string                   `json:"email" gorm:"-"`
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
	SaveDelivery bool                     `json:"save_delivery" gorm:"-"`
	Products     []OrderProductsCreateReq `json:"products" gorm:"-"`
	TimeDefault
}

type Orders struct {
	ID           int             `json:"id"`
	PlatformID   int             `json:"-"`
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

func (OrderCreateRequest) TableName() string {
	return "orders"
}

// 查詢功能
func (query *OrderQuery) GetCondition() *gorm.DB {
	sql := DB.Model(Orders{})
	sql.Where("member_id = ?", query.MemberID)

	return sql
}

// 基本CURD功能
func (req *OrderCreateRequest) Create() (err error) {
	req.Total = req.Price + req.Shipping
	err = DB.Create(&req).Error
	return
}

func (query *OrderQuery) FetchAll() (orders []Orders, pagination Pagination) {
	var count int64
	sql := query.GetCondition()
	sql.Count(&count)
	sql.Offset((query.Page - 1) * query.Items).Limit(query.Items).Order("created_at DESC").Scan(&orders)
	pagination = CreatePagination(query.Page, query.Items, count)
	return
}

// 關連功能
func (order *Orders) GetProducts() {
	DB.Model(&OrderProducts{}).Where("order_id = ?", order.ID).Order("created_at ASC").Scan(&order.Products)
}
