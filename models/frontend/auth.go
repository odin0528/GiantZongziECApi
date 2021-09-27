package frontend

import (
	. "eCommerce/internal/database"
	"os"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/plugin/soft_delete"
)

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type OAuthReq struct {
	Token    string `json:"token"`
	Platform string `json:"platform"`
}

type MemberToken struct {
	MemberID   int
	PlatformID int
	Token      string
	ExpiredAt  int
	CreatedAt  int
	DeletedAt  soft_delete.DeletedAt
}

type Claims struct {
	MemberID   int
	PlatformID int
	Nickname   string `json:"nickname"`
	jwt.StandardClaims
}

type FbPicture struct {
	Url string
}

type FbUserPicture struct {
	Data FbPicture
}

type FbUser struct {
	Name     string
	Email    string
	Birthday string
	ID       string
	Gender   string
	Picture  FbUserPicture
}

func (MemberToken) TableName() string {
	return "member_token"
}

func (token *MemberToken) CancelOldToken() {
	DB.Where("member_id = ?", token.MemberID).Delete(&MemberToken{})
}

func (token *MemberToken) Fetch() {
	DB.Model(MemberToken{}).Where("token = ? AND platform_id = ?", token.Token, token.PlatformID).Scan(token)
}

func GenerateToken(id int, platformID int, nickname string) (token string) {
	issuer := "GiantZongziEC"
	claims := Claims{
		MemberID:   id,
		PlatformID: platformID,
		Nickname:   nickname,
		StandardClaims: jwt.StandardClaims{
			Issuer: issuer,
		},
	}

	token, _ = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("JWT_SIGN")))
	return
}
