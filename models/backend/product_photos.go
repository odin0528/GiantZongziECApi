package backend

import (
	. "eCommerce/internal/database"
)

type ProductPhotos struct {
	ID         int    `json:"id"`
	CustomerID int    `json:"-"`
	ProductID  int    `json:"-"`
	Img        string `json:"img"`
	Sort       int    `json:"-"`
	TimeDefault
}

// 基本CURD功能
func (photo *ProductPhotos) Create() (err error) {
	err = DB.Create(&photo).Error
	return
}
