package backend

import (
	. "eCommerce/internal/database"
)

type ProductStyleTable struct {
	ID              int     `json:"id" gorm:"<-:create"`
	PlatformID      int     `json:"-" gorm:"<-:create"`
	ProductID       int     `json:"-" gorm:"<-:create"`
	GroupNo         int     `json:"-"`
	Sort            int     `json:"sort"`
	Title           string  `json:"title"`
	StyleTitle      string  `json:"style_title"`
	SubStyleTitle   string  `json:"sub_style_title"`
	Photo           string  `json:"photo"`
	Sku             string  `json:"sku"`
	Price           float64 `json:"price"`
	Qty             int     `json:"qty"`
	LowStock        int     `json:"low_stock"`
	Cost            float64 `json:"cost"`
	SuggestPrice    float64 `json:"suggest_price"`
	NoStoreDelivery int     `json:"no_store_delivery"`
	NoOverSale      bool    `json:"no_over_sale"`
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
	err = DB.
		Where("product_id = ? AND platform_id = ? AND id NOT IN (?)", style.ProductID, style.PlatformID, ids).
		Delete(&ProductStyleTable{}).Error
	return
}
