package backend

import (
	. "eCommerce/internal/database"
)

type Category struct {
	ID         int    `json:"-"`
	CustomerID int    `json:"-"`
	Title      string `json:"title"`
	TimeDefault
}

// 基本CURD功能
func (category *Category) Create() (err error) {
	category.ID = 0
	err = DB.Create(&category).Error
	return
}
