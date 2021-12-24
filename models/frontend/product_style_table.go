package frontend

import (
	. "eCommerce/internal/database"
)

type ProductStyleQuery struct {
	StyleID    int
	PlatformID int
	ProductID  int
}

type ProductStyleTable struct {
	ID              int     `json:"id" gorm:"<-:create"`
	PlatformID      int     `json:"-" gorm:"<-:create"`
	ProductID       int     `json:"-" gorm:"<-:create"`
	GroupNo         int     `json:"-"`
	Title           string  `json:"title"`
	StyleTitle      string  `json:"style_title"`
	SubStyleTitle   string  `json:"sub_style_title"`
	Photo           string  `json:"photo"`
	Sku             string  `json:"sku"`
	Price           float64 `json:"price"`
	Qty             int     `json:"qty"`
	WaitForDelivery int     `json:"-"`
	Cost            float64 `json:"cost"`
	SuggestPrice    float64 `json:"suggest_price"`
	NoStoreDelivery int     `json:"no_store_delivery"`
	NoOverSale      bool    `json:"no_over_sale"`
	Sold            int     `json:"sold"`
	TimeDefault
}

func (ProductStyleTable) TableName() string {
	return "product_style_table"
}

// 基本CURD功能
func (query *ProductStyleQuery) Fetch() (productStyleTable ProductStyleTable) {
	DB.Model(&ProductStyleTable{}).Where("platform_id = ? and product_id = ? and id = ?", query.PlatformID, query.ProductID, query.StyleID).
		First(&productStyleTable)
	return
}
