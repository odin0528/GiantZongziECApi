package backend

import (
	. "eCommerce/internal/database"
	"eCommerce/internal/rdb"
	"fmt"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type Brand struct {
	ID         int    `json:"id"`
	PlatformID int    `json:"-"`
	Title      string `json:"title"`
	LinkType   int    `json:"link_type"`
	Link       string `json:"link"`
	IsEnabled  bool   `json:"is_enabled"`
	Sort       int    `json:"sort"`
	DeletedAt  soft_delete.DeletedAt
	TimeDefault
}

type BrandQuery struct {
	ID         int
	PlatformID int
	SupplierID int `json:"supplier_id"`
	Pagination
}

func (Brand) TableName() string {
	return "brand"
}

func (query *BrandQuery) GetCondition() *gorm.DB {
	sql := DB.Model(Brand{})

	if query.ID != 0 {
		sql.Where("id = ?", query.ID)

	}
	if query.SupplierID != -1 {
		sql.Where("supplier_id = ?", query.SupplierID)
	}

	sql.Where("platform_id = ?", query.PlatformID)

	return sql
}

func (query *BrandQuery) Fetch() (brands []Brand, pagination Pagination) {
	var count int64
	sql := query.GetCondition()
	sql.Count(&count)
	sql.Offset((query.Page - 1) * query.Items).Limit(query.Items).Scan(&brands)
	pagination = CreatePagination(query.Page, query.Items, count)
	return
}

func (query *BrandQuery) FetchAll() (brands []Brand, err error) {

	key := fmt.Sprintf("platform_%d_brand", query.PlatformID)
	err = rdb.Get(key, &brands)
	if err == redis.Nil {
		sql := query.GetCondition()
		err = sql.Scan(&brands).Error

		if err != nil {
			return
		}

		rdb.Set(key, brands)
	}
	return
}