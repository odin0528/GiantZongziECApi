package backend

import (
	. "eCommerce/internal/database"

	"gorm.io/gorm"
)

type OrderListReq struct {
	PlatformID int `json:"-"`
	Status     int `json:"status"`
	Pagination
}

type OrderQuery struct {
	ID         int `json:"id"`
	PlatformID int `json:"-"`
	Status     int `json:"status"`
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

// 查詢功能
func (query *OrderQuery) GetCondition() *gorm.DB {
	sql := DB.Model(Orders{})
	if query.ID != 0 {
		sql.Where("id = ?", query.ID)
	}

	if query.Status != 0 {
		sql.Where("status = ?", query.Status)
	}

	sql.Where("platform_id = ?", query.PlatformID)

	return sql
}

func (query *OrderQuery) Fetch() (order Orders, err error) {
	sql := query.GetCondition()
	err = sql.First(&order).Error
	return
}

func (req *OrderListReq) FetchAll() (orders []Orders, pagination Pagination) {
	var count int64
	query := OrderQuery{
		PlatformID: req.PlatformID,
		Status:     req.Status,
	}
	sql := query.GetCondition()
	sql.Count(&count)
	sql.Offset((req.Page - 1) * req.Items).Limit(req.Items).Order("created_at DESC").Scan(&orders)
	pagination = CreatePagination(req.Page, req.Items, count)
	return
}

func (query *OrderQuery) FetchUntreated() (count int64) {
	sql := query.GetCondition()
	// 待付款，待出貨
	sql.Where("status IN (11, 21)")
	sql.Count(&count)
	return
}

// 關連功能
func (order *Orders) GetProducts() {
	DB.Model(&OrderProducts{}).Where("order_id = ?", order.ID).Order("created_at ASC").Scan(&order.Products)
}