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
	sql.Count(&count)
	sql.Offset((query.Page - 1) * query.Items).Limit(query.Items).Order("created_at DESC").Scan(&products)
	pagination = CreatePagination(query.Page, query.Items, count)
	return
}
