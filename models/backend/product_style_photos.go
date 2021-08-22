package backend

import (
	. "eCommerce/internal/database"
)

type ProductStylePhotos struct {
	ID             int    `json:"id"`
	PlatformID     int    `json:"-"`
	ProductID      int    `json:"-"`
	ProductStyleID int    `json:"product_style_id"`
	Img            string `json:"img"`
	TimeDefault
}

// 基本CURD功能
func (photo *ProductStylePhotos) Create() (err error) {
	err = DB.Create(&photo).Error
	return
}
