package frontend

import (
	. "eCommerce/internal/database"

	"gorm.io/plugin/soft_delete"
)

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type MemberToken struct {
	MemberID   int
	PlatformID int
	Token      string
	ExpiredAt  int
	CreatedAt  int
	DeletedAt  soft_delete.DeletedAt
}

func (MemberToken) TableName() string {
	return "member_token"
}

func (token *MemberToken) CancelOldToken() {
	DB.Debug().Where("member_id = ?", token.MemberID).Delete(&MemberToken{})
}

func (token *MemberToken) Fetch() {
	DB.Debug().Model(MemberToken{}).Where("token = ? AND platform_id = ?", token.Token, token.PlatformID).Scan(token)
}
