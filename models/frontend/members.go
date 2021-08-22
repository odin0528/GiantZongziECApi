package frontend

import (
	. "eCommerce/internal/database"
)

type MemberQuery struct {
	ID int
}

type Members struct {
	ID         int
	PlatformID int
	Email      string
	Password   string
	TimeDefault
}

func (query *MemberQuery) Fetch() (platform Platform) {
	DB.Model(&Platform{}).Select("id, title, logo_url, code").Where("hostname = ?", query.ID).Scan(&platform)
	return
}
