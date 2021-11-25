package frontend

import (
	. "eCommerce/internal/database"

	"gorm.io/plugin/soft_delete"
)

type Carts struct {
	ID         int `json:"-"`
	PlatformID int `json:"-"`
	MemberID   int `json:"-"`
	ProductID  int `json:"product_id"`
	StyleID    int `json:"style_id"`
	Qty        int `json:"qty"`
	DeletedAt  soft_delete.DeletedAt
	TimeDefault
}
type MemberCarts struct {
	ProductID     int    `json:"product_id"`
	StyleID       int    `json:"style_id"`
	Qty           int    `json:"qty"`
	Title         string `json:"title"`
	StyleTitle    string `json:"style_title"`
	SubStyleTitle string `json:"sub_style_title"`
	Photo         string `json:"photo"`
	Price         int    `json:"price"`
}

type CartsQuery struct {
	PlatformID int `json:"-"`
	MemberID   int `json:"-"`
}

func (Carts) TableName() string {
	return "carts"
}

func (query *CartsQuery) FetchAll() (carts []MemberCarts, err error) {
	err = DB.Debug().Table("carts").
		Select("carts.product_id, carts.style_id, carts.qty, product_style_table.title, product_style_table.style_title, product_style_table.sub_style_title, product_style_table.price, product_style_table.photo").
		Joins("inner join product_style_table on product_style_table.id = carts.style_id").
		Where("member_id = ? and carts.platform_id = ? AND carts.deleted_at = 0", query.MemberID, query.PlatformID).
		Order("carts.created_at ASC").
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
