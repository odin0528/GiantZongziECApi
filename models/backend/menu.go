package backend

import (
	// . "eCommerce/internal/database"

	"gorm.io/plugin/soft_delete"
)

type MenuMoveReq struct {
	ParentID  int `json:"parent_id"`
	Sort      int `json:"sort"`
	Direction int `json:"direction"`
}

type Menu struct {
	ID         int    `json:"id"`
	PlatformID int    `json:"-"`
	Title      string `json:"title"`
	LinkType   int    `json:"link_type"`
	Link       string `json:"link"`
	IsEnabled  bool   `json:"is_enabled"`
	Sort       int    `json:"sort"`
	DeletedAt  soft_delete.DeletedAt
	TimeDefault
}

func (Menu) TableName() string {
	return "platform_menu"
}
