package backend

import (
	"eCommerce/pkg/e"
	"encoding/base64"
	"net/http"
	"time"

	auth "eCommerce/internal/auth"
	. "eCommerce/internal/database"
	models "eCommerce/models/backend"

	"golang.org/x/crypto/scrypt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Login(c *gin.Context) {
	g := Gin{c}
	var req models.LoginReq
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	query := models.AdminQuery{
		Account: req.Account,
	}

	admin := query.Fetch()

	if admin.ID == 0 {
		g.Response(http.StatusOK, e.AccountNotExist, err)
		return
	}

	// 如果要重設密碼
	if admin.IsResetPwd {
		uuid := uuid.New()

		reset := models.AdminResetPassword{
			AdminID:   admin.ID,
			Token:     uuid.String(),
			ExpiredAt: int(time.Now().Unix()) + 60*10,
		}

		reset.CancelOldToken()
		DB.Create(&reset)

		g.Response(http.StatusOK, e.ResetRedirect, uuid.String())
		return
	}

	dk, err := scrypt.Key([]byte(req.Password), auth.Salt, 1<<15, 8, 1, 64)
	if base64.StdEncoding.EncodeToString(dk) != admin.Password {
		g.Response(http.StatusOK, e.AccountNotExist, nil)
		return
	}

	token := uuid.New().String()

	adminToken := models.AdminToken{
		AdminID: admin.ID,
		Token:   token,
	}
	adminToken.CancelOldToken()
	DB.Create(&adminToken)

	g.Response(http.StatusOK, e.Success, map[string]interface{}{"token": token, "title": admin.Title})
}

func ResetPassword(c *gin.Context) {
	g := Gin{c}
	var req models.ResetReq
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	if req.Password != req.CPassword {
		g.Response(http.StatusOK, e.PasswordNoMatch, err)
		return
	}

	reset := models.AdminResetPassword{
		Token: req.Token,
	}
	reset.Fetch()

	if reset.AdminID == 0 {
		g.Response(http.StatusOK, e.TokenNotExist, err)
		return
	}

	if time.Now().Unix() > int64(reset.ExpiredAt) {
		g.Response(http.StatusOK, e.TokenExpired, err)
		return
	}

	dk, err := scrypt.Key([]byte(req.Password), auth.Salt, 1<<15, 8, 1, 64)
	DB.Debug().Select("password", "is_reset_pwd").Where("id = ?", reset.AdminID).Updates(models.Admin{
		Password:   base64.StdEncoding.EncodeToString(dk),
		IsResetPwd: false,
	})

	reset.CancelOldToken()

	g.Response(http.StatusOK, e.Success, nil)
	return
}
