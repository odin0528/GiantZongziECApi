package frontend

import (
	"eCommerce/internal/auth"
	models "eCommerce/models/frontend"
	"eCommerce/pkg/e"
	"encoding/base64"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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

	var PlatformID int
	var MemberID int
	member := auth.JwtParser(c)
	if member == nil {
		pid, _ := c.Get("platform_id")
		PlatformID = pid.(int)
		MemberID = 0
	} else {
		PlatformID = member.PlatformID
		MemberID = member.MemberID
	}

	order.PlatformID = PlatformID
	order.MemberID = MemberID

	if len(order.Email) > 0 && MemberID == 0 {

		dk, _ := scrypt.Key([]byte(order.Phone), auth.Salt, 1<<15, 8, 1, 64)

		member := models.Members{
			Email:      order.Email,
			Password:   base64.StdEncoding.EncodeToString(dk),
			PlatformID: PlatformID,
		}

		err = DB.Create(&member).Error
		if err != nil {
			g.Response(http.StatusOK, e.EmailDuplicate, err)
			return
		}
		order.MemberID = member.ID

		// 一併幫會員做登入
		issuer := "GiantZongziEC"
		claims := models.Claims{
			MemberID:   member.ID,
			PlatformID: PlatformID,
			Nickname:   member.Nickname,
			StandardClaims: jwt.StandardClaims{
				Issuer: issuer,
			},
		}

		token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(os.Getenv("JWT_SIGN")))

	} else if MemberID != 0 {
		carts := models.Carts{
			MemberID:   MemberID,
			PlatformID: PlatformID,
		}
		carts.Clean()

		if order.SaveDelivery {
			delivery := models.MemberDelivery{
				PlatformID:   PlatformID,
				MemberID:     MemberID,
				Fullname:     order.Fullname,
				Phone:        order.Phone,
				Address:      order.Address,
				Memo:         order.Memo,
				Method:       order.Method,
				StoreID:      order.StoreID,
				StoreName:    order.StoreName,
				StoreAddress: order.StoreAddress,
				StorePhone:   order.StorePhone,
			}

			DB.Create(&delivery)
		}
	}

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
				StyleTitle: style.StyleTitle,
				Photo:      style.Photo,
				Sku:        style.Sku,
			}

			orderProduct.Create()
		}
	}

	g.Response(http.StatusOK, e.Success, token)
}
