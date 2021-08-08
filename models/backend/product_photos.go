package backend

import (
	. "eCommerce/internal/database"
)

type ProductPhotos struct {
	ID         int    `json:"id"`
	CustomerID int    `json:"-" gorm:"<-:create"`
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
	err = DB.Debug().Save(&photo).Error
	return
}

func (photo *ProductPhotos) Delete() (err error) {
	err = DB.Debug().Where("id = ? AND customer_id = ?", photo.ID, photo.CustomerID).Delete(&ProductPhotos{}).Error
	return
}

func (photo *ProductPhotos) Fetch() (err error) {
	err = DB.Debug().Table("product_photos").Where("id = ? AND customer_id = ?", photo.ID, photo.CustomerID).Scan(&photo).Error
	return
}
