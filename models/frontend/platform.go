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
}

func (Platform) TableName() string {
	return "platform"
}

func (query *PlatformQuery) Fetch() (platform Platform) {
	DB.Model(&Platform{}).Select("id, title, logo_url, code").Where("hostname = ?", query.Hostname).Scan(&platform)
	return
}
