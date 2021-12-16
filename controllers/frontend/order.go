package frontend

import (
	"context"
	"eCommerce/internal/auth"
	"eCommerce/internal/line"
	models "eCommerce/models/frontend"
	"eCommerce/pkg/e"
	"encoding/base64"
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/Laysi/go-ecpay-sdk"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gotokatsuya/line-pay-sdk-go/linepay"
	"golang.org/x/crypto/scrypt"
	"gorm.io/gorm"

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

	tx := DB.Begin()

	errCode := OrderValidation(tx, PlatformID, &order)
	if errCode != 200 {
		g.Response(http.StatusOK, errCode, order.Products)
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

		err = tx.Create(&member).Error
		if err != nil {
			tx.Rollback()
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
			County:       order.County,
			District:     order.District,
			ZipCode:      order.ZipCode,
			Address:      order.Address,
			Memo:         order.Memo,
			Method:       order.Method,
			StoreID:      order.StoreID,
			StoreName:    order.StoreName,
			StoreAddress: order.StoreAddress,
			StorePhone:   order.StorePhone,
		}

		err = tx.Create(&delivery).Error
		if err != nil {
			tx.Rollback()
			g.Response(http.StatusInternalServerError, e.StatusInternalServerError, err)
			return
		}
	}

	switch order.Payment {
	case 2:
		order.Status = 21
	default:
		order.Status = 11
	}

	// 新增訂單
	order.Total = order.Price + order.Shipping - order.Discount
	err = tx.Create(&order).Error
	if err != nil {
		tx.Rollback()
		g.Response(http.StatusInternalServerError, e.StatusInternalServerError, err)
		return
	}

	for pi, product := range order.Products {
		for si, style := range product.Styles {
			itemName = append(itemName, fmt.Sprintf("%s %s", product.Title, style.StyleTitle))
			order.Products[pi].Styles[si].Title = product.Title
			orderProduct := models.OrderProducts{
				OrderID:         order.ID,
				ProductID:       product.ProductID,
				StyleID:         style.StyleID,
				Qty:             style.BuyCount,
				Price:           style.Price,
				IsDiscount:      style.IsDiscount,
				Discount:        style.Discount,
				DiscountedPrice: style.DiscountedPrice,
				Total:           float64(style.BuyCount) * style.DiscountedPrice,
				Title:           product.Title,
				StyleTitle:      style.StyleTitle,
				Photo:           style.Photo,
				Sku:             style.Sku,
			}

			err = tx.Create(&orderProduct).Error

			if err.Error() == "out_of_stock" {
				tx.Rollback()
				g.Response(http.StatusOK, e.UpdateFailForOutOfStock, err)
				return
			}

			if err != nil {
				tx.Rollback()
				g.Response(http.StatusInternalServerError, e.StatusInternalServerError, err)
				return
			}
		}
	}

	if order.Payment == 4 {

		platform, _ := c.Get("platform")
		Platform := platform.(models.Platform)
		orderUuid := uuid.New().String()

		var pay *linepay.Client

		if os.Getenv("ENV") != "production" {
			pay, _ = linepay.New(os.Getenv("LINE_PAY_ID"), os.Getenv("LINE_PAY_KEY"), linepay.WithSandbox())
		} else {
			pay, _ = linepay.New(os.Getenv("LINE_PAY_ID"), os.Getenv("LINE_PAY_KEY"))
		}

		requestReq := &linepay.RequestRequest{
			Amount:   int(order.Total),
			Currency: "TWD",
			OrderID:  orderUuid,
			Packages: []*linepay.RequestPackage{},
			RedirectURLs: &linepay.RequestRedirectURLs{
				ConfirmURL: fmt.Sprintf(os.Getenv("EC_CHECKOUT_FINISH_URL"), c.Request.Header["Hostname"][0]),
				CancelURL:  fmt.Sprintf(os.Getenv("EC_CHECKOUT_CANCEL_URL"), c.Request.Header["Hostname"][0]),
			},
		}

		for _, product := range order.Products {
			for _, style := range product.Styles {
				requestReq.Packages = append(requestReq.Packages,
					&linepay.RequestPackage{
						ID:     fmt.Sprintf("%d", order.ID),
						Amount: style.BuyCount * int(style.DiscountedPrice),
						Name:   Platform.Title,
						Products: []*linepay.RequestPackageProduct{
							&linepay.RequestPackageProduct{
								ID:       fmt.Sprintf("%d-%d", product.ProductID, style.StyleID),
								Name:     fmt.Sprintf("%s %s", product.Title, style.StyleTitle),
								Quantity: style.BuyCount,
								Price:    int(style.DiscountedPrice),
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

		if !order.IsFreeShipping {
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
		}

		requestResp, _, _ := pay.Request(context.Background(), requestReq)

		if requestResp.ReturnCode == "0000" {
			DB.Model(models.Orders{}).Where("id = ?", order.ID).
				Updates(map[string]interface{}{"order_uuid": orderUuid, "transaction_id": requestResp.Info.TransactionID})
		}

		g.Response(http.StatusOK, e.Success, map[string]interface{}{"token": token, "payment": requestResp.Info.PaymentURL, "request": requestReq})
		return
	} else if order.Payment == 2 {
		// 貨到付款就等後台出託運單
		carts := models.Carts{
			MemberID:   MemberID,
			PlatformID: PlatformID,
		}
		carts.Clean()

		// email.SendOrderNotify(order)
		line.SendOrderNotifyByOrderCreateRequest(order)

		g.Response(http.StatusOK, e.Success, map[string]interface{}{"token": token})
	} else {
		var client *ecpay.Client
		p, _ := c.Get("platform")
		platform := p.(models.Platform)
		if os.Getenv("ENV") != "production" {
			client = ecpay.NewStageClient(
				ecpay.WithReturnURL(fmt.Sprintf("%s%s", os.Getenv("API_URL"), os.Getenv("ECPAY_PAYMENT_FINISH_URL"))),
				ecpay.WithClientBackURL(fmt.Sprintf(os.Getenv("EC_URL"), c.Request.Header["Hostname"][0])),
				ecpay.WithOrderResultURL(fmt.Sprintf(os.Getenv("ECPAY_CLIENT_RETURN_URL"), c.Request.Header["Hostname"][0], "%2Fcheckout%2Ffinish")),
				ecpay.WithDebug,
			)
		} else {
			client = ecpay.NewClient(
				os.Getenv("ECPAY_MERCHANT_ID"),
				os.Getenv("ECPAY_MERCHANT_HASH_KEY"),
				os.Getenv("ECPAY_MERCHANT_HASH_IV"),
				fmt.Sprintf("%s%s", os.Getenv("API_URL"), os.Getenv("ECPAY_PAYMENT_FINISH_URL")),
				ecpay.WithClientBackURL(fmt.Sprintf(os.Getenv("EC_URL"), c.Request.Header["Hostname"][0])),
				ecpay.WithOrderResultURL(fmt.Sprintf(os.Getenv("ECPAY_CLIENT_RETURN_URL"), c.Request.Header["Hostname"][0], "%2Fcheckout%2Ffinish")),
				ecpay.WithDebug,
			)
		}
		aio := client.CreateOrder(fmt.Sprintf("%s%d", os.Getenv("ECPAY_MERCHANT_TRADE_NO_PREFIX"), order.ID), time.Now(), int(order.Total), fmt.Sprintf("%s %s", platform.Title, order.Memo), itemName)

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
		html, _ := aio.GenerateRequestHtml()

		carts := models.Carts{
			MemberID:   MemberID,
			PlatformID: PlatformID,
		}
		carts.Clean()

		g.Response(http.StatusOK, e.Success, map[string]interface{}{"token": token, "ecpay": html})
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

	order.Status = req.Status
	DB.Select("status").Save(&order)

	g.Response(http.StatusOK, e.Success, order)
}

func OrderValidation(tx *gorm.DB, PlatformID int, order *models.OrderCreateRequest) int {
	priceChange, outOfStock := false, false
	count := 0
	var total, shipping, checkoutPercent, checkoutDiscount, productDiscount, shippingDiscount float64 = 0, 0, 100, 0, 0, 0
	isFreeShipping := false

	// 檢查價格或庫存
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

			// 判斷價格是否有異動或正確
			if item.NoOverSale && item.Qty < style.BuyCount {
				order.Products[productIndex].Styles[styleIndex].Qty = item.Qty
				outOfStock = true
			}

			count += style.BuyCount
			total += float64(style.BuyCount) * style.Price
		}
	}

	if outOfStock {
		// return e.OutOfStock
	}

	if priceChange || total != order.Price {
		return e.ProductPriceChange
	}

	// 檢查有沒有調整運費
	platform := &models.Platform{
		ID: PlatformID,
	}
	logistics := platform.GetLogistics()
	switch order.Method {
	case 1:
		shipping = logistics.HomeChargeFee
	case 2:
		shipping = logistics.UniChargeFee
	case 3:
		shipping = logistics.FamilyChargeFee
	case 4:
		shipping = logistics.HilifeChargeFee
	case 5:
		shipping = logistics.OkChargeFee
	}

	if shipping != order.Shipping {
		return e.ShippingChange
	}

	order.Qty = count

	// 取得優惠活動，並算完折扣
	promotions := models.GetPromotionByID(PlatformID)

	// 先算產品折扣
	for _, promotion := range promotions {
		switch promotion.Type {
		case "special_price":
			checkProductDiscount := CheckSpecialPrice(promotion, order.Products)
			if checkProductDiscount == -1 {
				return e.PromotionChange
			} else {
				productDiscount += checkProductDiscount
			}
		}
	}

	// 後算購物車折扣
	for _, promotion := range promotions {
		switch promotion.Type {
		case "sitewide_discount":
			if promotion.Mode == "total_qty" && promotion.Qty <= count {
				if promotion.Method == "percent" {
					checkoutPercent *= promotion.Percent / 100
				} else if promotion.Method == "discount" {
					checkoutDiscount += promotion.Discount
				}
			} else if promotion.Mode == "total_price" && promotion.Money <= total-productDiscount {
				if promotion.Method == "percent" {
					checkoutPercent *= promotion.Percent / 100
				} else if promotion.Method == "discount" {
					checkoutDiscount += promotion.Discount
				}
			}
		case "free_shipping":
			if promotion.Mode == "total_qty" && promotion.Qty <= count {
				isFreeShipping = true
			} else if promotion.Mode == "total_price" && promotion.Money <= total-productDiscount {
				isFreeShipping = true
			}
		}
	}

	if isFreeShipping {
		shippingDiscount = shipping
	}
	order.IsFreeShipping = isFreeShipping

	if total-math.Round((total-productDiscount)*(checkoutPercent/100)-checkoutDiscount)+shippingDiscount != order.Discount {
		return e.PromotionChange
	}

	return e.Success
}

func CheckSpecialPrice(promotion models.Promotions, products []models.OrderProductsCreateReq) float64 {
	var discount float64
	discount = 0
	for _, product := range products {
		for _, style := range product.Styles {
			discountedPrice := math.Round(style.Price * promotion.Percent / 100)
			if style.DiscountedPrice != discountedPrice || style.Discount != style.Price-discountedPrice {
				return -1
			}
			discount += style.Discount * float64(style.BuyCount)
		}
	}

	return discount
}
