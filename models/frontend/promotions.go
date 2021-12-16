package frontend

import (
	. "eCommerce/internal/database"
)

type PromotionListReq struct {
	PlatformID int `json:"-"`
	Pagination
}

type Promotions struct {
	ID             int     `json:"id"`
	PlatformID     int     `json:"-" gorm:"<-:create"`
	Title          string  `json:"title"`
	StartTimestamp int     `json:"start_timestamp"`
	EndTimestamp   int     `json:"end_timestamp"`
	Type           string  `json:"type"`
	Mode           string  `json:"mode"`
	Method         string  `json:"method"`
	Qty            int     `json:"qty"`
	Money          float64 `json:"money"`
	Percent        float64 `json:"percent"`
	Discount       float64 `json:"discount"`
	IsEnabled      bool    `json:"is_enabled"`
	DeletedAt      int     `json:"-"`
	TimeDefault
}

func GetPromotionByID(ID int) (promotions []Promotions) {
	DB.Model(&Promotions{}).
		Where("platform_id = ? AND is_enabled = 1 AND start_timestamp <= UNIX_TIMESTAMP() AND end_timestamp > UNIX_TIMESTAMP()", ID).
		Scan(&promotions)
	return
}
