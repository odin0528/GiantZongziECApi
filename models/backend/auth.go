package backend

import (
	. "eCommerce/internal/database"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/plugin/soft_delete"
)

type LoginReq struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type ResetReq struct {
	Password  string `json:"password"`
	CPassword string `json:"cPassword"`
	Token     string `json:"token"`
}

type AdminResetPassword struct {
	AdminID   int
	Token     string
	ExpiredAt int
	CreatedAt int
	DeletedAt soft_delete.DeletedAt
}

type AdminToken struct {
	AdminID   int
	Token     string
	ExpiredAt int
	CreatedAt int
	DeletedAt soft_delete.DeletedAt
}

type Claims struct {
	AdminID    int
	PlatformID int
	Title      string `json:"title"`
	jwt.StandardClaims
}

func (AdminResetPassword) TableName() string {
	return "admin_reset_password"
}
func (AdminToken) TableName() string {
	return "admin_token"
}

func (reset *AdminResetPassword) CancelOldToken() {
	DB.Debug().Where("admin_id = ?", reset.AdminID).Delete(&AdminResetPassword{})
}

func (token *AdminToken) CancelOldToken() {
	DB.Debug().Where("admin_id = ?", token.AdminID).Delete(&AdminToken{})
}

func (reset *AdminResetPassword) Fetch() {
	DB.Debug().Model(AdminResetPassword{}).Where("token = ?", reset.Token).Scan(reset)
}

func (token *AdminToken) Fetch() {
	DB.Debug().Model(AdminToken{}).Where("token = ?", token.Token).Scan(token)
}
