package backend

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
func (photo *ProductPhotos) Create() (err error) {
	err = DB.Create(&photo).Error
	return
}

func (photo *ProductPhotos) Update() (err error) {
	err = DB.Save(&photo).Error
	return
}

func (photo *ProductPhotos) Delete() (err error) {
	err = DB.Where("id = ? AND platform_id = ?", photo.ID, photo.PlatformID).Delete(&ProductPhotos{}).Error
	return
}

func (photo *ProductPhotos) Fetch() (err error) {
	err = DB.Table("product_photos").Where("id = ? AND platform_id = ?", photo.ID, photo.PlatformID).Scan(&photo).Error
	return
}
