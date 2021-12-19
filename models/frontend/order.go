package frontend

import (
	. "eCommerce/internal/database"

	"gorm.io/gorm"
)

type OrderQuery struct {
	MemberID      int    `json:"-"`
	TransactionID string `json:"-"`
	OrderUuid     string `json:"-"`
	Status        int    `json:"-"`
	Pagination
}

type OrderCreateRequest struct {
	ID             int                      `json:"-"`
	PlatformID     int                      `json:"-"`
	MemberID       int                      `json:"-"`
	Email          string                   `json:"email" gorm:"-"`
	Fullname       string                   `json:"fullname"`
	Phone          string                   `json:"phone"`
	County         string                   `json:"county"`
	District       string                   `json:"district"`
	ZipCode        string                   `json:"zip_code"`
	Address        string                   `json:"address"`
	Memo           string                   `json:"memo"`
	Method         int                      `json:"method"`
	Total          float64                  `json:"-"`
	Price          float64                  `json:"price"`
	Discount       float64                  `json:"discount"`
	Shipping       float64                  `json:"shipping"`
	IsFreeShipping bool                     `json:"-"`
	Qty            int                      `json:"-"`
	Payment        int                      `json:"payment"`
	StoreID        string                   `json:"store_id"`
	StoreName      string                   `json:"store_name"`
	StoreAddress   string                   `json:"store_address"`
	StorePhone     string                   `json:"store_phone"`
	Status         int                      `json:"-"`
	SaveDelivery   bool                     `json:"save_delivery" gorm:"-"`
	Products       []OrderProductsCreateReq `json:"products" gorm:"-"`
	TimeDefault
}

type Orders struct {
	ID              int             `json:"id"`
	PlatformID      int             `json:"-"`
	MemberID        int             `json:"-"`
	Fullname        string          `json:"fullname"`
	Phone           string          `json:"phone"`
	County          string          `json:"county"`
	District        string          `json:"district"`
	ZipCode         string          `json:"zip_code"`
	Address         string          `json:"address"`
	Memo            string          `json:"memo"`
	Method          int             `json:"method"`
	Total           float64         `json:"total"`
	Price           float64         `json:"price"`
	Discount        float64         `json:"discount"`
	Shipping        float64         `json:"shipping"`
	Qty             int             `json:"qty"`
	Payment         int             `json:"payment"`
	StoreID         string          `json:"store_id"`
	StoreName       string          `json:"store_name"`
	StoreAddress    string          `json:"store_address"`
	StorePhone      string          `json:"store_phone"`
	Status          int             `json:"status"`
	ShipmentNo      string          `json:"shipment_no"`
	LogisticsStatus int             `json:"logistics_status"`
	LogisticsMsg    string          `json:"logistics_msg"`
	Products        []OrderProducts `json:"products" gorm:"-"`
	TimeDefault
}

type OrderUpdateReq struct {
	TransactionID string `json:"transaction_id"`
	OrderUuid     string `json:"order_id"`
	Status        int    `json:"status"`
}

func (OrderCreateRequest) TableName() string {
	return "orders"
}

// 查詢功能
func (query *OrderQuery) GetCondition() *gorm.DB {
	sql := DB.Model(Orders{})

	if query.MemberID != 0 {
		sql.Where("member_id = ?", query.MemberID)
	}

	if query.TransactionID != "" {
		sql.Where("transaction_id = ?", query.TransactionID)
	}

	if query.OrderUuid != "" {
		sql.Where("order_uuid = ?", query.OrderUuid)
	}

	if query.Status != 0 {
		sql.Where("status = ?", query.Status)
	}

	return sql
}

// 基本CURD功能
func (req *OrderCreateRequest) Create() (err error) {
	req.Total = req.Price + req.Shipping - req.Discount
	err = DB.Create(&req).Error
	return
}

func (query *OrderQuery) Fetch() (order Orders) {
	query.GetCondition().Scan(&order)
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
