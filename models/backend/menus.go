package backend

import (
	. "eCommerce/internal/database"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type MenuMoveReq struct {
	Sort      int `json:"sort"`
	Direction int `json:"direction"`
}

type Menus struct {
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

func (Menus) TableName() string {
	return "platform_menu"
}

type MenuQuery struct {
	ID         int
	PlatformID int
	Sort       int
}

func (menus *Menus) Delete() (err error) {
	DB.Delete(menus)
	err = DB.Table("platform_menu").Where("platform_id = ? AND sort > ?", menus.PlatformID, menus.Sort).Update("sort", gorm.Expr("sort - 1")).Error
	return
}

func (query *MenuQuery) GetCondition() *gorm.DB {
	sql := DB.Model(Menus{})

	if query.ID != 0 {
		sql.Where("id = ?", query.ID)
	}

	if query.Sort != 0 {
		sql.Where("sort = ?", query.Sort)
	}

	sql.Where("platform_id = ?", query.PlatformID)

	return sql
}

func (query *MenuQuery) Fetch() (menu Menus) {
	sql := query.GetCondition()
	sql.First(&menu)
	return
}

func (query *MenuQuery) FetchAll() (menus []Menus, err error) {
	sql := query.GetCondition().Order("sort ASC")
	err = sql.Scan(&menus).Error
	return
}
