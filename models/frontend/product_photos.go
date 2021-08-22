package frontend

import (
	. "eCommerce/internal/database"
)

type ProductPhotos struct {
	ID         int    `json:"id"`
	PlatformID int    `json:"-" gorm:"<-:create"`
	ProductID  int    `json:"-" gorm:"<-:create"`
	Img        string `json:"photo"`
	Sort       int    `json:"-"`
	TimeDefault
}

// 基本CURD功能
func (photo *ProductPhotos) Fetch() (err error) {
	err = DB.Debug().Table("product_photos").Where("id = ? AND platform_id = ?", photo.ID, photo.PlatformID).Scan(&photo).Error
	return
}
