package backend

import (
	. "eCommerce/internal/database"
)

type ReportType struct {
	Type       string `uri:"type"`
	PlatformID int    `json:"-"`
	Pagination
}

func (query *ReportType) LowStock() (products []Products, pagination Pagination) {
	var count int64
	sql := DB.Debug().Table("report_low_stock")
	sql.Where("platform_id = ?", query.PlatformID)
	sql.Count(&count)
	sql.Offset((query.Page - 1) * query.Items).Limit(query.Items).Order("created_at DESC").Scan(&products)
	pagination = CreatePagination(query.Page, query.Items, count)
	return
}

func (query *ReportType) OverSale() (products []Products, pagination Pagination) {
	var count int64
	sql := DB.Debug().Table("report_over_sale")
	sql.Where("platform_id = ?", query.PlatformID)
	sql.Count(&count)
	sql.Offset((query.Page - 1) * query.Items).Limit(query.Items).Order("created_at DESC").Scan(&products)
	pagination = CreatePagination(query.Page, query.Items, count)
	return
}
