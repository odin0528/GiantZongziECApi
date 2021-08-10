package frontend

import (
	. "eCommerce/internal/database"
	"eCommerce/pkg/e"
)

type PageReq struct {
	Page       string `json:"page" uri:"page"`
	CustomerID int    `json:"customer_id"`
}

type Pages struct {
	PageID     int  `json:"page_id"`
	CustomerID int  `json:"-"`
	ReleasedAt int  `json:"released_at"`
	IsEnabled  bool `json:"is_enabled"`
	TimeDefault
}

func (req *PageReq) Fetch() (pages Pages) {
	DB.Table("rel_customer_pages").Where("page_id = ? and customer_id = ?", e.PageList[req.Page], req.CustomerID).Scan(&pages)
	return
}

func (pages *Pages) Validate() bool {
	// data is not exist
	if pages.PageID == 0 || !pages.IsEnabled {
		return false
	}
	return true
}
