package backend

import (
	. "eCommerce/internal/database"
)

type ProductSubStyle struct {
	ID         int    `json:"id" gorm:"<-:create"`
	PlatformID int    `json:"-" gorm:"<-:create"`
	ProductID  int    `json:"-" gorm:"<-:create"`
	Title      string `json:"title"`
	Sort       int    `json:"sort"`
	TimeDefault
}

func (ProductSubStyle) TableName() string {
	return "product_sub_style"
}

// 基本CURD功能
func (style *ProductSubStyle) Create() (err error) {
	style.ID = 0
	err = DB.Create(&style).Error
	return
}

func (style *ProductSubStyle) Update() (err error) {
	err = DB.Debug().Save(&style).Error
	return
}

func (style *ProductSubStyle) DeleteNotExistStyle(ids []int) (err error) {
	sql := DB.Debug().Where("product_id = ? AND platform_id = ?", style.ProductID, style.PlatformID)
	if len(ids) > 0 {
		sql.Where("id NOT IN (?)", ids)
	}
	err = sql.Delete(&ProductSubStyle{}).Error
	return
}
