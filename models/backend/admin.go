package backend

import (
	. "eCommerce/internal/database"

	"gorm.io/gorm"
)

type AdminQuery struct {
	ID         int
	PlatformID int
	Account    string
	Title      string
}

type Admin struct {
	ID         int
	PlatformID int
	Account    string
	Password   string
	Title      string
	IsResetPwd bool
	TimeDefault
}

func (Admin) TableName() string {
	return "admin"
}

func (query *AdminQuery) Query() *gorm.DB {
	sql := DB.Table("admin")
	if query.ID != 0 {
		sql.Where("id = ?", query.ID)
	}

	if query.PlatformID != 0 {
		sql.Where("platform_id = ?", query.PlatformID)
	}

	if query.Account != "" {
		sql.Where("account like ?", query.Account)
	}

	if query.Title != "" {
		sql.Where("title like '%?%'", query.Title)
	}

	return sql
}

func (query *AdminQuery) Fetch() (admin Admin) {
	sql := query.Query()
	sql.First(&admin)
	return
}
