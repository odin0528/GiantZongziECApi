package backend

import (
	"context"
	"eCommerce/pkg/e"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"

	models "eCommerce/models/backend"

	. "eCommerce/internal/database"
	"eCommerce/internal/line"

	"github.com/gin-gonic/gin"
	"github.com/gotokatsuya/line-pay-sdk-go/linepay"
)

func LinePayFinish(c *gin.Context) {
	g := Gin{c}
	var req *models.OrderLinepayReq
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	query := models.OrderQuery{
		TransactionID: req.TransactionID,
		OrderUuid:     req.OrderUuid,
		Status:        11,
	}

	order, _ := query.FetchLinePayOrder()

	if order.ID == 0 {
		g.Response(http.StatusBadRequest, e.OrderNotExist, err)
		return
	}

	transactionID, err := strconv.ParseInt(req.TransactionID, 10, 64)

	requestReq := &linepay.CheckPaymentStatusRequest{}
	var pay *linepay.Client

	if os.Getenv("ENV") != "production" {
		pay, _ = linepay.New(os.Getenv("LINE_PAY_ID"), os.Getenv("LINE_PAY_KEY"), linepay.WithSandbox())
	} else {
		pay, _ = linepay.New(os.Getenv("LINE_PAY_ID"), os.Getenv("LINE_PAY_KEY"))
	}

	res, _, _ := pay.CheckPaymentStatus(context.Background(), transactionID, requestReq)

	if res.ReturnCode == "0110" {
		confirmReq := &linepay.ConfirmRequest{
			Amount:   int(order.Total),
			Currency: "TWD",
		}

		confirmRes, _, _ := pay.Confirm(context.Background(), transactionID, confirmReq)

		if confirmRes.ReturnCode == "0000" {
			order.Status = 21
			order.PaymentChargeFee = math.Floor(order.Total * 0.0315)
			DB.Select("status", "payment_charge_fee").Updates(&order)

			order.GetProducts()

			line.SendOrderNotifyByOrder(order)

			// 第三方付款成功後，清掉會員的購物車
			if order.MemberID != 0 {
				carts := models.Carts{
					MemberID:   order.MemberID,
					PlatformID: order.PlatformID,
				}
				carts.Clean()
			}

			g.Response(http.StatusOK, e.Success, nil)
		} else {
			g.Response(http.StatusOK, e.StatusInternalServerError, fmt.Sprintf("Line PAY付款確認失敗，錯誤代碼：%s，錯誤訊息：%s", confirmRes.ReturnCode, confirmRes.ReturnMessage))
		}
	} else {
		g.Response(http.StatusOK, e.StatusInternalServerError, fmt.Sprintf("Line PAY付款確認失敗，錯誤代碼：%s，錯誤訊息：%s", res.ReturnCode, res.ReturnMessage))
	}
}
