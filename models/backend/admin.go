package backend

import (
	. "eCommerce/internal/database"

	"gorm.io/gorm"
)

type AdminQuery struct {
	ID         int
	CustomerID int
	Account    string
	Title      string
}

type Admin struct {
	ID         int
	CustomerID int
	Account    string
	Password   string `json:"-" gorm:"<-:create"`
	Title      string
	IsResetPwd bool
	TimeDefault
}

func (query *AdminQuery) Query() *gorm.DB {
	sql := DB.Table("admin")
	if query.ID != 0 {
		sql.Where("id = ?", query.ID)
	}

	if query.CustomerID != 0 {
		sql.Where("customer_id = ?", query.CustomerID)
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
