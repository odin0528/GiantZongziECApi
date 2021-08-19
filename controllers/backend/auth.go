package backend

import (
	"eCommerce/pkg/e"
	"net/http"
	"time"

	. "eCommerce/internal/database"
	models "eCommerce/models/backend"

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
		g.Response(http.StatusOK, e.StatusNotFound, err)
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

	g.Response(http.StatusOK, e.Success, admin)
}
