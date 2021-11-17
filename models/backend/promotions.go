package backend

import (
	. "eCommerce/internal/database"
)

type PromotionListReq struct {
	PlatformID int `json:"-"`
	Pagination
}

type Promotions struct {
	ID             int    `json:"id"`
	PlatformID     int    `json:"-" gorm:"<-:create"`
	Title          string `json:"title"`
	StartTimestamp int    `json:"start_timestamp"`
	EndTimestamp   int    `json:"end_timestamp"`
	Type           string `json:"type"`
	Mode           string `json:"mode"`
	Method         string `json:"method"`
	Qty            int    `json:"qty"`
	Money          int    `json:"money"`
	Percent        int    `json:"percent"`
	Discount       int    `json:"discount"`
	IsEnabled      bool   `json:"is_enabled"`
	DeletedAt      int    `json:"-"`
	TimeDefault
}

func (req *PromotionListReq) FetchAll() (promotions []Promotions, pagination Pagination) {
	var count int64
	sql := DB.Model(&Promotions{}).Where("platform_id = ?", req.PlatformID)
	sql.Count(&count)
	sql.Offset((req.Page - 1) * req.Items).Limit(req.Items).Scan(&promotions)
	pagination = CreatePagination(req.Page, req.Items, count)
	return
}
