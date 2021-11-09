package backend

import (
	"eCommerce/pkg/e"
	"fmt"
	"net/http"
	"os"
	"strings"

	. "eCommerce/internal/database"
	"eCommerce/internal/money"
	models "eCommerce/models/backend"

	"github.com/gin-gonic/gin"
	"github.com/liudng/godump"
)

func OrderFetch(c *gin.Context) {
	g := Gin{c}
	var query models.ProductQuery
	err := c.ShouldBindUri(&query)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}
	PlatformID, _ := c.Get("platform_id")
	query.PlatformID = PlatformID.(int)
	product := query.Fetch()
	product.GetPhotos()
	product.GetStyle()
	product.GetSubStyle()
	product.GetStyleTable()

	g.Response(http.StatusOK, e.Success, product)
}

func OrderList(c *gin.Context) {
	g := Gin{c}
	var req models.OrderListReq
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}
	PlatformID, _ := c.Get("platform_id")
	req.PlatformID = PlatformID.(int)
	orders, pagination := req.FetchAll()

	for index := range orders {
		orders[index].GetProducts()
	}

	g.PaginationResponse(http.StatusOK, e.Success, orders, pagination)
}

func OrderNextStep(c *gin.Context) {
	g := Gin{c}
	var query models.OrderQuery
	err := c.BindJSON(&query)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}
	PlatformID, _ := c.Get("platform_id")
	query.PlatformID = PlatformID.(int)
	order, err := query.Fetch()

	if err != nil {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}

	switch order.Status {
	case 11:
		order.Status = 21
	case 21:
		order.Status = 31
	case 31:
		order.Status = 99
	}

	err = DB.Select("status").Updates(&order).Error

	if err != nil {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}

	g.Response(http.StatusOK, e.Success, nil)
}

func OrderMakeShipmentNo(c *gin.Context) {
	g := Gin{c}
	var query models.OrderQuery
	err := c.BindJSON(&query)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}
	PlatformID, _ := c.Get("platform_id")
	query.PlatformID = PlatformID.(int)
	query.Status = 21
	order, err := query.Fetch()

	if err != nil {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}

	order.GetProducts()

	response, err := money.CreateLogisticsOrder(order)
	if err != nil {
		g.Response(http.StatusOK, e.StatusInternalServerError, err.Error())
		return
	}

	order.LogisticsID = response.Get("AllPayLogisticsID")
	order.LogisticsStatus = 1
	order.LogisticsMsg = response.Get("RtnMsg")
	err = DB.Select("status", "logistics_id", "logistics_status", "logistics_msg").Updates(&order).Error

	if err != nil {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}

	g.Response(http.StatusOK, e.Success, nil)
}

func OrderShipmentPrint(c *gin.Context) {
	g := Gin{c}
	var query models.BatchOrderQuery
	err := c.BindJSON(&query)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}
	PlatformID, _ := c.Get("platform_id")
	query.PlatformID = PlatformID.(int)
	query.Status = 22
	orders, err := query.FetchAll()

	if err != nil {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}

	info, _ := money.QueryLogisticsInfoV2(orders[0].LogisticsID)
	godump.Dump(info)
	return

	ids := []string{}
	tradeNo := []string{}
	paymentNo := []string{}
	validationNo := []string{}
	shipmentNo := []string{}

	for _, order := range orders {
		ids = append(ids, order.LogisticsID)
		tradeNo = append(tradeNo, fmt.Sprintf("%s%d", os.Getenv("ECPAY_MERCHANT_TRADE_NO_PREFIX"), order.ID))
		shipmentNo = append(shipmentNo, order.ShipmentNo)
		paymentNo = append(paymentNo, order.ShipmentNo[:len(order.ShipmentNo)-4])
		validationNo = append(validationNo, order.ShipmentNo[len(order.ShipmentNo)-4:])
	}

	params := map[string]string{}
	var url string

	switch orders[0].Method {
	case 1:
		params["MerchantTradeNo"] = ""
		url = "https://logistics.ecpay.com.tw/helper/PrintTradeDocument"
	case 2:
		params["CVSPaymentNo"] = strings.Join(paymentNo, ",")
		params["CVSValidationNo"] = strings.Join(validationNo, ",")
		url = "https://logistics.ecpay.com.tw/Express/PrintUniMartC2COrderInfo"
	case 3:
		params["CVSPaymentNo"] = strings.Join(paymentNo, ",")
		params["CVSValidationNo"] = strings.Join(validationNo, ",")
		url = "https://logistics.ecpay.com.tw/helper/PrintTradeDocument"
	case 4:
		params["MerchantTradeNo"] = strings.Join(tradeNo, ",")
		url = "https://logistics.ecpay.com.tw/helper/PrintTradeDocument"
	case 5:
		params["CVSPaymentNo"] = strings.Join(paymentNo, ",")
		params["CVSValidationNo"] = strings.Join(validationNo, ",")
		url = "https://logistics.ecpay.com.tw/Express/PrintOKMARTC2COrderInfo"
	}

	params["AllPayLogisticsID"] = strings.Join(ids, ",")
	params["MerchantID"] = os.Getenv("ECPAY_MERCHANT_ID")
	params["isPreview"] = "True"

	checkMac := money.MakeLogisticsCheckMac(params)
	params["CheckMacValue"] = checkMac

	formBody := ""
	for k, v := range params {
		formBody += fmt.Sprintf(`<input type="hidden" name="%s" value="%s" />`, k, v)
	}

	formString := `<form id="PostForm" name="PostForm" action="%s" method="POST" target="_blank">%s</form>`
	g.Response(http.StatusOK, e.Success, fmt.Sprintf(formString, url, formBody))
}

func OrderUntreated(c *gin.Context) {
	g := Gin{c}
	var query models.OrderQuery
	PlatformID, _ := c.Get("platform_id")
	query.PlatformID = PlatformID.(int)
	count := query.FetchUntreated()

	g.Response(http.StatusOK, e.Success, count)
}
