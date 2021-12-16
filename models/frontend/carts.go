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
	BuyCount   int `json:"buy_count"`
	DeletedAt  soft_delete.DeletedAt
	TimeDefault
}
type MemberCarts struct {
	ProductID       int    `json:"product_id"`
	StyleID         int    `json:"style_id"`
	BuyCount        int    `json:"buy_count"`
	Qty             int    `json:"qty"`
	Title           string `json:"title"`
	StyleTitle      string `json:"style_title"`
	SubStyleTitle   string `json:"sub_style_title"`
	NoOverSale      bool   `json:"no_over_sale"`
	NoStoreDelivery int    `json:"no_store_delivery"`
	Photo           string `json:"photo"`
	Price           int    `json:"price"`
}

type CartsQuery struct {
	PlatformID int `json:"-"`
	MemberID   int `json:"-"`
}

type GuestCartsQuery struct {
	StyleID    []int `json:"ids"`
	PlatformID int   `json:"-"`
}

func (Carts) TableName() string {
	return "carts"
}

func (query *CartsQuery) FetchAll() (carts []MemberCarts, err error) {
	err = DB.Table("carts").
		Select("carts.product_id, carts.style_id, carts.buy_count, pst.title, pst.style_title, pst.sub_style_title, pst.price, pst.photo, pst.qty, pst.no_over_sale, pst.no_store_delivery").
		Joins("inner join product_style_table as pst on pst.id = carts.style_id").
		Where("member_id = ? and carts.platform_id = ? AND carts.deleted_at = 0", query.MemberID, query.PlatformID).
		Order("carts.created_at ASC").
		Scan(&carts).Error

	return
}

func (query *GuestCartsQuery) FetchAll() (carts []MemberCarts, err error) {
	err = DB.Table("product_style_table").
		Select("id as style_id, product_id, title, style_title, sub_style_title, price, photo, qty, no_over_sale, no_store_delivery").
		Where("id IN ? AND platform_id = ?", query.StyleID, query.PlatformID).
		Scan(&carts).Error

	return
}

func (req *Carts) Update() error {
	return DB.Select("buy_count").
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
