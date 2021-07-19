package backend

import (
	. "eCommerce/internal/database"
)

type PageReq struct {
	CustomerID int `json:"customer_id"`
	Pagination
}

type Pages struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (req *PageReq) GetPageList() (pagesRowset []Pages, err error) {
	err = DB.Debug().Table("rel_customer_pages as rel").Select("pages.*").Joins("inner join pages on rel.page_id = pages.id").
		Where("rel.customer_id = ?", req.CustomerID).
		Scan(&pagesRowset).Error
	return
}
