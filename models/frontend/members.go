package frontend

import (
	. "eCommerce/internal/database"

	"gorm.io/gorm"
)

type MemberQuery struct {
	ID         int
	PlatformID int
	Email      string
}

type Members struct {
	ID         int    `json:"-" gorm:"<-create"`
	PlatformID int    `json:"-"`
	Email      string `json:"email"`
	Password   string `json:"-"`
	Nickname   string `json:"nickname"`
	Phone      string `json:"phone"`
	Birthday   string `json:"birthday"`
	TimeDefault
}

func (query *MemberQuery) GetCondition() *gorm.DB {
	sql := DB.Debug().Model(Members{})

	if query.Email != "" {
		sql.Where("email like ?", query.Email)
	}

	sql.Where("platform_id = ?", query.PlatformID)

	return sql
}

func (query *MemberQuery) Fetch() (member Members) {
	sql := query.GetCondition()
	sql.Scan(&member)
	return
}
