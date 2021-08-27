package backend

import (
	. "eCommerce/internal/database"

	"gorm.io/gorm"
)

type MemberQuery struct {
	PlatformID int
	Pagination
}

type Members struct {
	ID         int    `json:"id"`
	PlatformID int    `json:"-"`
	Email      string `json:"email"`
	Nickname   string `json:"nickname"`
	Phone      string `json:"phone"`
	Birthday   string `json:"birthday"`
	TimeDefault
}

func (Members) TableName() string {
	return "members"
}

// 查詢功能
func (query *MemberQuery) GetCondition() *gorm.DB {
	sql := DB.Model(Members{})

	sql.Where("platform_id = ?", query.PlatformID)

	return sql
}

func (query *MemberQuery) Fetch() (order Orders, err error) {
	sql := query.GetCondition()
	err = sql.First(&order).Error
	return
}

func (query *MemberQuery) FetchAll() (members []Members, pagination Pagination) {
	var count int64
	sql := query.GetCondition()
	sql.Count(&count)
	sql.Offset((query.Page - 1) * query.Items).Limit(query.Items).Order("created_at DESC").Scan(&members)
	pagination = CreatePagination(query.Page, query.Items, count)
	return
}
