package frontend

import (
	. "eCommerce/internal/database"
	"time"

	"gorm.io/gorm"
)

type MemberQuery struct {
	ID         int
	PlatformID int
	Email      string
}

type MemberRegisterReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Phone    string `json:"phone"`
	Birthday string `json:"birthday"`
}

type Members struct {
	ID         int       `json:"-" gorm:"<-create"`
	PlatformID int       `json:"-"`
	Email      string    `json:"email"`
	Password   string    `json:"-"`
	Nickname   string    `json:"nickname" gorm:"default:null"`
	Phone      string    `json:"phone" gorm:"default:null"`
	Birthday   time.Time `json:"birthday" gorm:"default:null"`
	TimeDefault
}

func (query *MemberQuery) GetCondition() *gorm.DB {
	sql := DB.Model(Members{})

	if query.ID != 0 {
		sql.Where("id = ?", query.ID)
	}

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

func (req *MemberRegisterReq) Create() {
	DB.Model(Members{}).Create(&req)
}
