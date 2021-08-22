package frontend

import (
	"eCommerce/internal/auth"
	models "eCommerce/models/frontend"
	"eCommerce/pkg/e"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/scrypt"

	. "eCommerce/internal/database"
)

func OrderCreate(c *gin.Context) {
	g := Gin{c}
	var order *models.OrderCreateRequest
	var token string
	err := c.BindJSON(&order)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	PlatformID, _ := c.Get("platform_id")
	MemberID, _ := c.Get("member_id")

	if len(order.Email) > 0 && MemberID == 0 {

		dk, _ := scrypt.Key([]byte(order.Phone), auth.Salt, 1<<15, 8, 1, 64)

		member := models.Members{
			Email:      order.Email,
			Password:   base64.StdEncoding.EncodeToString(dk),
			PlatformID: PlatformID.(int),
		}

		err = DB.Create(&member).Error
		if err != nil {
			g.Response(http.StatusOK, e.EmailDuplicate, err)
			return
		}
		order.MemberID = member.ID

		// 一併幫會員做登入
		token = uuid.New().String()

		memberToken := models.MemberToken{
			MemberID: member.ID,
			Token:    token,
		}
		memberToken.CancelOldToken()
		DB.Create(&memberToken)

	}

	order.PlatformID = PlatformID.(int)
	switch order.Payment {
	case 1:
		order.Status = 11
	case 2:
		order.Status = 21
	}
	order.Create()

	for _, product := range order.Products {
		for _, style := range product.Styles {
			orderProduct := models.OrderProducts{
				OrderID:    order.ID,
				ProductID:  product.ProductID,
				StyleID:    style.StyleID,
				Qty:        style.Qty,
				Price:      style.Price,
				Total:      float32(style.Qty) * style.Price,
				Title:      product.Title,
				StyleTitle: style.Title + style.SubTitle,
				Photo:      style.Photo,
				Sku:        style.Sku,
			}

			orderProduct.Create()
		}
	}

	g.Response(http.StatusOK, e.Success, token)
}
