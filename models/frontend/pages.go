package frontend

import (
	. "eCommerce/internal/database"

	"gorm.io/plugin/soft_delete"
)

type PageReq struct {
	Url        string `json:"url" uri:"page"`
	PlatformID int    `json:"platform_id"`
}

type Pages struct {
	ID         int    `json:"id"`
	PlatformID int    `json:"-"`
	Type       int    `json:"type"`
	Url        string `json:"url"  validate:"required,alphanum"`
	Title      string `json:"title"  validate:"required"`
	IsMenu     bool   `json:"is_menu"`
	IsEnabled  bool   `json:"is_enabled"`
	ReleasedAt int    `json:"released_at"`
	DeletedAt  soft_delete.DeletedAt
	TimeDefault
}

func (req *PageReq) Fetch() (pages Pages, err error) {
	err = DB.Model(&Pages{}).Where("url = ? and platform_id = ?", req.Url, req.PlatformID).Scan(&pages).Error
	return
}
