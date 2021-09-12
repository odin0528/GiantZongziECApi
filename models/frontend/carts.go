package frontend

import (
	. "eCommerce/internal/database"

	"gorm.io/plugin/soft_delete"
)

type Carts struct {
	ID         int     `json:"-"`
	PlatformID int     `json:"-"`
	MemberID   int     `json:"-"`
	ProductID  int     `json:"product_id"`
	StyleID    int     `json:"style_id"`
	Qty        int     `json:"qty"`
	Price      float32 `json:"price"`
	Total      float32 `json:"total"`
	Title      string  `json:"title"`
	StyleTitle string  `json:"style_title"`
	Photo      string  `json:"photo"`
	Sku        string  `json:"sku"`
	DeletedAt  soft_delete.DeletedAt
	TimeDefault
}

type CartsQuery struct {
	PlatformID int `json:"-"`
	MemberID   int `json:"-"`
}

func (Carts) TableName() string {
	return "carts"
}

func (query *CartsQuery) FetchAll() (carts []Carts, err error) {
	err = DB.Model(&Carts{}).
		Where("member_id = ? and platform_id = ?", query.MemberID, query.PlatformID).
		Order("created_at DESC").
		Scan(&carts).Error

	return
}

func (req *Carts) Update() error {
	return DB.Select("qty").
		Where("member_id = ? and platform_id = ? and product_id = ? and style_id = ?", req.MemberID, req.PlatformID, req.ProductID, req.StyleID).
		Updates(&req).Error
}

func (req *Carts) Delete() error {
	return DB.
		Where("member_id = ? and platform_id = ? and product_id = ? and style_id = ?", req.MemberID, req.PlatformID, req.ProductID, req.StyleID).
		Delete(&Carts{}).Error
}

func (req *Carts) Clean() error {
	return DB.
		Where("member_id = ? and platform_id = ?", req.MemberID, req.PlatformID).
		Delete(&Carts{}).Error
}
