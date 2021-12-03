package backend

import (
	. "eCommerce/internal/database"
	"fmt"

	"gorm.io/gorm"
)

/*
LogisticsStatus：
1xx:正流程
2xx:逆流程
100: 賣家已出貨
110: 物流配送中
120: 商品已到店，待取貨
130: 買家已取貨

200: 買家未取貨，退貨中
210: 物流退貨中
220: 已退回寄貨地
230: 賣家已取貨

Status：
11: 待付款
21: 待出貨
22: 揀貨中
23: 已出貨
51: 配送中
61: 退貨中

91: 訂單完成
92: 退貨完成 (買家未取件)

98: 訂單取消 (超過7天未付款)
99: 訂單取消 (消費者取消付款)
*/

type Orders struct {
	ID               int             `json:"id"`
	PlatformID       int             `json:"-"`
	MemberID         int             `json:"-"`
	Fullname         string          `json:"fullname"`
	Phone            string          `json:"phone"`
	County           string          `json:"county"`
	District         string          `json:"district"`
	ZipCode          string          `json:"zip_code"`
	Address          string          `json:"address"`
	Memo             string          `json:"memo"`
	Method           int             `json:"method"`
	Total            float64         `json:"total"`
	Price            float64         `json:"price"`
	Shipping         float64         `json:"shipping"`
	Discount         float64         `json:"discount"`
	Qty              int             `json:"qty"`
	IsFreeShipping   bool            `json:"is_free_shipping"`
	Payment          int             `json:"payment"`
	PaymentChargeFee float64         `json:"-"`
	StoreID          string          `json:"store_id"`
	StoreName        string          `json:"store_name"`
	StoreAddress     string          `json:"store_address"`
	StorePhone       string          `json:"store_phone"`
	Status           int             `json:"status"`
	LogisticsID      string          `json:"logistics_id"`
	ShipmentNo       string          `json:"shipment_no"`
	LogisticsStatus  int             `json:"logistics_status"`
	LogisticsMsg     string          `json:"logistics_msg"`
	Products         []OrderProducts `json:"products" gorm:"-"`
	TimeDefault
}

type OrderListReq struct {
	IDs             []int  `json:"ids"`
	PlatformID      int    `json:"-"`
	Status          int    `json:"status"`
	LogisticsStatus int    `json:"logistics_status"`
	Method          int    `json:"method"`
	WithoutProducts bool   `json:"without_products"`
	OrderBy         string `json:"order_by"`
	Sort            string `json:"sort"`
	Pagination
}

type OrderQuery struct {
	ID            int    `json:"id"`
	PlatformID    int    `json:"-"`
	TransactionID string `json:"-"`
	OrderUuid     string `json:"-"`
	LogisticsID   string `json:"-"`
	Status        int    `json:"status"`
}

type BatchOrderQuery struct {
	IDs             []int `json:"ids"`
	PlatformID      int   `json:"-"`
	Status          int   `json:"status"`
	LogisticsStatus int   `json:"logistics_status"`
}

type OrderLinepayReq struct {
	OrderUuid     string `json:"order_id"`
	TransactionID string `json:"transaction_id"`
	Status        int    `json:"status"`
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

func (query *BatchOrderQuery) GetCondition() *gorm.DB {
	sql := DB.Model(Orders{})
	if len(query.IDs) > 0 {
		sql.Where("id IN ?", query.IDs)
	}

	if query.Status != 0 {
		sql.Where("status = ?", query.Status)
	}

	sql.Where("platform_id = ?", query.PlatformID)

	return sql
}

func (query *OrderListReq) GetCondition() *gorm.DB {
	sql := DB.Debug().Model(Orders{})

	if len(query.IDs) > 0 {
		sql.Where("id IN ?", query.IDs)
	}

	if query.LogisticsStatus != 0 {
		sql.Where("logistics_status = ?", query.LogisticsStatus)
	}

	if query.Status != 0 {
		sql.Where("status = ?", query.Status)
	}

	if query.Method != 0 {
		sql.Where("method = ?", query.Method)
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
	sql := req.GetCondition()
	sql.Count(&count)
	sql.Offset((req.Page - 1) * req.Items).Limit(req.Items).Order(fmt.Sprintf("%s %s", req.OrderBy, req.Sort)).Scan(&orders)
	pagination = CreatePagination(req.Page, req.Items, count)
	return
}

func (query *BatchOrderQuery) FetchAll() (orders []Orders, err error) {
	sql := query.GetCondition()
	err = sql.Scan(&orders).Error
	return
}

func (query *OrderQuery) FetchLinePayOrder() (order Orders, err error) {
	sql := DB.Model(Orders{}).Where("transaction_id = ? AND order_uuid = ? AND status = ?", query.TransactionID, query.OrderUuid, query.Status)
	err = sql.First(&order).Error
	return
}

func (query *OrderQuery) FetchForLogistics() (order Orders, err error) {
	sql := DB.Model(Orders{}).Where("id = ?", query.ID)
	err = sql.First(&order).Error
	return
}

func (query *OrderQuery) FetchUntreated() (count int64) {
	sql := query.GetCondition()
	// 待付款，待出貨
	sql.Where("status = 21")
	sql.Count(&count)
	return
}

// 關連功能
func (order *Orders) GetProducts() {
	DB.Model(&OrderProducts{}).Where("order_id = ?", order.ID).Order("created_at ASC").Scan(&order.Products)
}
