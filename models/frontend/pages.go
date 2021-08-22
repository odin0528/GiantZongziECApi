package frontend

import (
	. "eCommerce/internal/database"
	"eCommerce/pkg/e"
)

type PageReq struct {
	Page       string `json:"page" uri:"page"`
	PlatformID int    `json:"platform_id"`
}

type Pages struct {
	PageID     int  `json:"page_id"`
	PlatformID int  `json:"-"`
	ReleasedAt int  `json:"released_at"`
	IsEnabled  bool `json:"is_enabled"`
	TimeDefault
}

func (req *PageReq) Fetch() (pages Pages) {
	DB.Table("rel_platform_pages").Where("page_id = ? and platform_id = ?", e.PageList[req.Page], req.PlatformID).Scan(&pages)
	return
}

func (pages *Pages) Validate() bool {
	// data is not exist
	if pages.PageID == 0 || !pages.IsEnabled {
		return false
	}
	return true
}
