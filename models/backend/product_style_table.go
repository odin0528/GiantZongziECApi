package backend

import (
	. "eCommerce/internal/database"
)

type ProductStyleTable struct {
	ID         int    `json:"id" gorm:"<-:create"`
	CustomerID int    `json:"-" gorm:"<-:create"`
	ProductID  int    `json:"-" gorm:"<-:create"`
	Group      int    `json:"-"`
	Title      string `json:"title"`
	SubTitle   string `json:"subTitle"`
	Sku        string `json:"sku"`
	Price      int    `json:"price"`
	Qty        int    `json:"qty"`
	TimeDefault
}

func (ProductStyleTable) TableName() string {
	return "product_style_table"
}

// 基本CURD功能
func (style *ProductStyleTable) Create() (err error) {
	err = DB.Create(&style).Error
	return
}

func (style *ProductStyleTable) Update() (err error) {
	err = DB.Save(&style).Error
	return
}

func (style *ProductStyleTable) DeleteNotExistStyle(ids []int) (err error) {
	err = DB.Debug().
		Where("product_id = ? AND customer_id = ? AND id NOT IN (?)", style.ProductID, style.CustomerID, ids).
		Delete(&ProductStyleTable{}).Error
	return
}
