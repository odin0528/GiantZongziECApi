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
	ID         int     `json:"id" gorm:"<-:create"`
	PlatformID int     `json:"-" gorm:"<-:create"`
	ProductID  int     `json:"-" gorm:"<-:create"`
	Group      int     `json:"-"`
	Title      string  `json:"title"`
	SubTitle   string  `json:"subTitle"`
	Sku        string  `json:"sku"`
	Price      float32 `json:"price"`
	Qty        int     `json:"qty"`
	TimeDefault
}

func (ProductStyleTable) TableName() string {
	return "product_style_table"
}

// 基本CURD功能
func (query *ProductStyleQuery) Fetch() (productStyleTable ProductStyleTable) {
	DB.Debug().Model(&ProductStyleTable{}).Where("platform_id = ? and product_id = ? and id = ?", query.PlatformID, query.ProductID, query.StyleID).
		First(&productStyleTable)
	return
}
