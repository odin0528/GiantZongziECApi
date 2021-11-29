package backend

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

type CartsQuery struct {
	PlatformID int `json:"-"`
	MemberID   int `json:"-"`
}

func (Carts) TableName() string {
	return "carts"
}

func (req *Carts) Clean() error {
	return DB.
		Where("member_id = ? and platform_id = ?", req.MemberID, req.PlatformID).
		Delete(&Carts{}).Error
}
