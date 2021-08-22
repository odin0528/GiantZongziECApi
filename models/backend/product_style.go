package backend

import (
	. "eCommerce/internal/database"
)

type ProductStyle struct {
	ID         int    `json:"id" gorm:"<-:create"`
	PlatformID int    `json:"-" gorm:"<-:create"`
	ProductID  int    `json:"-" gorm:"<-:create"`
	Title      string `json:"title"`
	Img        string `json:"photo"`
	Sort       int    `json:"sort"`
	TimeDefault
}

func (ProductStyle) TableName() string {
	return "product_style"
}

// 基本CURD功能
func (style *ProductStyle) Create() (err error) {
	style.ID = 0
	err = DB.Create(&style).Error
	return
}

func (style *ProductStyle) Update() (err error) {
	err = DB.Debug().Save(&style).Error
	return
}

func (style *ProductStyle) DeleteNotExistStyle(ids []int) (err error) {
	sql := DB.Debug().Where("product_id = ? AND platform_id = ?", style.ProductID, style.PlatformID)
	if len(ids) > 0 {
		sql.Where("id NOT IN (?)", ids)
	}
	err = sql.Delete(&ProductStyle{}).Error
	return
}
