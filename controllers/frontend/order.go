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
	var order models.OrderCreateRequest
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

	errCode := OrderValidation(PlatformID, &order)
	if errCode != 200 {
		if errCode == e.ProductPriceChange {
			if MemberID != 0 {
				for _, product := range order.Products {
					for _, style := range product.Styles {
						carts := models.Carts{
							MemberID:   MemberID,
							PlatformID: PlatformID,
							ProductID:  product.ProductID,
							StyleID:    style.StyleID,
							Price:      style.Price,
						}
						DB.Debug().Model(&carts).Where("platform_id = ? and member_id = ? and product_id = ? and style_id = ? and deleted_at = 0",
							PlatformID, MemberID, product.ProductID, style.StyleID).
							Update("price", carts.Price)
					}
				}
			}
			g.Response(http.StatusOK, errCode, order.Products)
		}
		return
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

func OrderValidation(PlatformID int, order *models.OrderCreateRequest) int {
	priceChange := false
	count := 0
	var total, shipping, percent, discount float32 = 0, 0, 100, 0
	for productIndex, product := range order.Products {
		for styleIndex, style := range product.Styles {
			query := &models.ProductStyleQuery{
				StyleID:    style.StyleID,
				PlatformID: PlatformID,
				ProductID:  product.ProductID,
			}
			item := query.Fetch()

			// 判斷價格是否有異動或正確
			if item.Price != style.Price {
				order.Products[productIndex].Styles[styleIndex].Price = item.Price
				priceChange = true
			}

			count += style.Qty
			total += float32(style.Qty) * style.Price
		}
	}

	if priceChange {
		return e.ProductPriceChange
	}

	order.Qty = count

	switch order.Method {
	case 1:
		shipping = 120
	default:
		shipping = 60
	}

	if total != order.Price || shipping != order.Shipping {
		return e.ProductPriceChange
	}

	// 取得優惠活動，並算完折扣
	promotions := models.GetPromotionByID(PlatformID)
	for _, promotion := range promotions {
		switch promotion.Type {
		case "sitewide_discount":
			if promotion.Mode == "total_qty" && promotion.Qty <= count {
				if promotion.Method == "percent" {
					percent *= promotion.Percent / 100
				} else if promotion.Method == "discount" {
					discount += promotion.Discount
				}
			} else if promotion.Mode == "total_price" && promotion.Money <= total {
				if promotion.Method == "percent" {
					percent *= promotion.Percent / 100
				} else if promotion.Method == "discount" {
					discount += promotion.Discount
				}
			}
		}
	}

	if total-(total*(percent/100)-discount) != order.Discount {
		return e.PromotionChange
	}

	return e.Success
}
