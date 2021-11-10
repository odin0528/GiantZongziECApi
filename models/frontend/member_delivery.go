package frontend

import (
	. "eCommerce/internal/database"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type MemberDeliveryQuery struct {
	MemberID int `json:"-"`
}

type MemberDelivery struct {
	ID           int    `json:"id"`
	PlatformID   int    `json:"-"`
	MemberID     int    `json:"-"`
	Fullname     string `json:"fullname"`
	Phone        string `json:"phone"`
	County       string `json:"county"`
	District     string `json:"district"`
	ZipCode      string `json:"zip_code"`
	Address      string `json:"address"`
	Memo         string `json:"memo"`
	Method       int    `json:"method"`
	StoreID      string `json:"store_id"`
	StoreName    string `json:"store_name"`
	StoreAddress string `json:"store_address"`
	StorePhone   string `json:"store_phone"`
	DeletedAt    soft_delete.DeletedAt
	TimeDefault
}

func (MemberDelivery) TableName() string {
	return "member_delivery"
}

// 基本CURD功能
func (req *MemberDelivery) Create() (err error) {
	err = DB.Create(&req).Error
	return
}

// 查詢功能
func (query *MemberDeliveryQuery) GetCondition() *gorm.DB {
	sql := DB.Model(MemberDelivery{})
	sql.Where("member_id = ?", query.MemberID)

	return sql
}

func (query *MemberDeliveryQuery) FetchAll() (deliveries []MemberDelivery) {
	sql := query.GetCondition()
	sql.Scan(&deliveries)
	return
}
