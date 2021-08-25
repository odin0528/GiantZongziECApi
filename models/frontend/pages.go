package frontend

import (
	. "eCommerce/internal/database"

	"gorm.io/plugin/soft_delete"
)

type PageReq struct {
	Url        int `json:"url" uri:"url"`
	PlatformID int `json:"platform_id"`
}

type Pages struct {
	ID         int    `json:"page_id"`
	PlatformID int    `json:"-"`
	Url        string `json:"url"`
	Title      string `json:"title"`
	IsMenu     bool   `json:"is_menu"`
	ReleasedAt int    `json:"released_at"`
	DeletedAt  soft_delete.DeletedAt
	TimeDefault
}

func (req *PageReq) Fetch() (pages Pages, err error) {
	err = DB.Debug().Model(&Pages{}).Where("url = ? and platform_id = ?", req.Url, req.PlatformID).Scan(&pages).Error
	return
}
