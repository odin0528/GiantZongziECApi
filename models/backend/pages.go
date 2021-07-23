package backend

import (
	. "eCommerce/internal/database"
)

type PageReq struct {
	CustomerID int `json:"customer_id"`
	Pagination
}

type Pages struct {
	PageID     int    `json:"page_id"`
	Name       string `json:"name"`
	ReleasedAt int    `json:"released_at"`
	TimeDefault
}

func (req *PageReq) GetPageList() (pagesRowset []Pages, err error) {
	err = DB.Table("rel_customer_pages as rel").Select("rel.*, pages.name").Joins("inner join pages on rel.page_id = pages.id").
		Where("rel.customer_id = ?", req.CustomerID).
		Scan(&pagesRowset).Error
	return
}
