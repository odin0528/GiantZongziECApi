package frontend

import (
	"context"
	"eCommerce/internal/auth"
	models "eCommerce/models/frontend"
	"eCommerce/pkg/e"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Laysi/go-ecpay-sdk"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gotokatsuya/line-pay-sdk-go/linepay"
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
	var itemName []string
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
		token = models.GenerateToken(member)

	} else if MemberID != 0 && order.SaveDelivery {
		// 如果本身有登入，又有勾儲存運送方式
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

	switch order.Payment {
	case 2:
		order.Status = 21
	default:
		order.Status = 11
	}
	order.Create()

	for _, product := range order.Products {
		for _, style := range product.Styles {
			itemName = append(itemName, fmt.Sprintf("%s %s", product.Title, style.StyleTitle))
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

	if order.Payment == 4 {

		platform, _ := c.Get("platform")
		Platform := platform.(models.Platform)
		orderUuid := uuid.New().String()

		pay, _ := linepay.New(os.Getenv("LINE_PAY_ID"), os.Getenv("LINE_PAY_KEY"), linepay.WithSandbox())
		// pay, _ := linepay.New(os.Getenv("LINE_PAY_ID"), os.Getenv("LINE_PAY_KEY"))
		requestReq := &linepay.RequestRequest{
			Amount:   int(order.Price + order.Shipping),
			Currency: "TWD",
			OrderID:  orderUuid,
			Packages: []*linepay.RequestPackage{},
			RedirectURLs: &linepay.RequestRedirectURLs{
				ConfirmURL: fmt.Sprintf(os.Getenv("LINE_PAYMENT_FINISH_URL"), c.Request.Header["Hostname"][0]),
				CancelURL:  fmt.Sprintf(os.Getenv("LINE_PAYMENT_CANCEL_URL"), c.Request.Header["Hostname"][0]),
			},
		}

		for _, product := range order.Products {
			for _, style := range product.Styles {
				requestReq.Packages = append(requestReq.Packages,
					&linepay.RequestPackage{
						ID:     fmt.Sprintf("%d", order.ID),
						Amount: style.Qty * int(style.Price),
						Name:   Platform.Title,
						Products: []*linepay.RequestPackageProduct{
							&linepay.RequestPackageProduct{
								ID:       fmt.Sprintf("%d-%d", product.ProductID, style.StyleID),
								Name:     fmt.Sprintf("%s %s", product.Title, style.StyleTitle),
								Quantity: style.Qty,
								Price:    int(style.Price),
								ImageURL: style.Photo,
							},
						},
					},
				)
			}
		}

		/* requestReq.Packages[0].Products = append(requestReq.Packages[0].Products, &linepay.RequestPackageProduct{
			ID:       fmt.Sprintf("order-shipping-%d", order.ID),
			Name:     "運費",
			Quantity: 1,
			Price:    int(order.Shipping),
		}) */

		requestReq.Packages = append(requestReq.Packages,
			&linepay.RequestPackage{
				ID:     fmt.Sprintf("%d", order.ID),
				Amount: int(order.Shipping),
				Name:   Platform.Title,
				Products: []*linepay.RequestPackageProduct{
					&linepay.RequestPackageProduct{
						ID:       fmt.Sprintf("%d-%s", order.ID, "shipping"),
						Name:     "運費",
						Quantity: 1,
						Price:    int(order.Shipping),
					},
				},
			},
		)

		requestResp, _, _ := pay.Request(context.Background(), requestReq)

		if requestResp.ReturnCode == "0000" {
			DB.Debug().Model(models.Orders{}).Where("id = ?", order.ID).
				Updates(map[string]interface{}{"order_uuid": orderUuid, "transaction_id": requestResp.Info.TransactionID})
		}

		g.Response(http.StatusOK, e.Success, map[string]interface{}{"token": token, "payment": requestResp.Info.PaymentURL, "request": requestReq})
		return
	} else if order.Payment == 2 {
		// 如果不是三方支付，交易完就先清
		carts := models.Carts{
			MemberID:   MemberID,
			PlatformID: PlatformID,
		}
		carts.Clean()
		g.Response(http.StatusOK, e.Success, map[string]interface{}{"token": token})
	} else {
		p, _ := c.Get("platform")
		platform := p.(models.Platform)
		client := ecpay.NewStageClient(
			ecpay.WithReturnURL(fmt.Sprintf("%s%s", os.Getenv("API_URL"), os.Getenv("ECPAY_PAYMENT_FINISH_URL"))),
			ecpay.WithOrderResultURL(fmt.Sprintf(os.Getenv("ECPAY_CLIENT_RETURN_URL"), c.Request.Header["Hostname"][0], "%2Fcheckout%2Ffinish")),
			ecpay.WithDebug,
		)
		aio := client.CreateOrder(fmt.Sprintf("GZEC%d", order.ID), time.Now(), int(order.Total), fmt.Sprintf("%s %s", platform.Title, order.Memo), itemName)

		switch order.Payment {
		case 3:
			aio.SetCreditPayment()
		case 1:
			aio.SetWebAtmPayment()
		case 5:
			aio.SetAtmPayment()
		case 6:
			aio.SetCvsPayment()
		case 7:
			aio.SetBarcodePayment()
		}
		mac, _ := aio.GenerateCheckMac()
		html, _ := aio.GenerateRequestHtml()

		DB.Model(&order).Update("ecpay_mac", mac)

		g.Response(http.StatusOK, e.Success, map[string]interface{}{"token": token, "ecpay": html, "mac": mac})
	}

}

func OrderUpdate(c *gin.Context) {
	g := Gin{c}
	var req *models.OrderUpdateReq
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	query := models.OrderQuery{
		TransactionID: req.TransactionID,
		OrderUuid:     req.OrderUuid,
	}

	order := query.Fetch()

	if order.ID == 0 {
		g.Response(http.StatusBadRequest, e.OrderNotExist, err)
		return
	}

	if order.Status != 11 {
		g.Response(http.StatusOK, e.OrderIsFinish, err)
		return
	}

	switch req.Status {
	case 21:
		// 第三方付款成功後，清掉會員的購物車
		if order.MemberID != 0 {
			carts := models.Carts{
				MemberID:   order.MemberID,
				PlatformID: order.PlatformID,
			}
			carts.Clean()
		}
	}

	order.Status = req.Status
	DB.Select("status").Save(&order)

	g.Response(http.StatusOK, e.Success, order)
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
