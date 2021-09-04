package frontend

import (
	. "eCommerce/internal/database"
)

type PlatformQuery struct {
	Hostname string
}

type Platform struct {
	ID      int    `json:"-"`
	Title   string `json:"title"`
	LogoUrl string `json:"logo_url"`
	Code    string `json:"code"`
	FBPixel string `json:"fb_pixel"`
}

func (Platform) TableName() string {
	return "platform"
}

func (query *PlatformQuery) Fetch() (platform Platform) {
	DB.Model(&Platform{}).Select("id, title, logo_url, code").Where("hostname = ?", query.Hostname).Scan(&platform)
	return
}

func (platform *Platform) GetMenu() (pages []Pages) {
	DB.Model(&Pages{}).Where("platform_id = ? AND is_menu = 1 AND is_enabled = 1 AND released_at > 0", platform.ID).Scan(&pages)
	return
}
