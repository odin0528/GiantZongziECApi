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

	if req.WithoutProducts != true {
		for index := range orders {
			orders[index].GetProducts()
		}
	}

	g.PaginationResponse(http.StatusOK, e.Success, orders, pagination)
}

func OrderNextStep(c *gin.Context) {
	g := Gin{c}
	var query models.BatchOrderQuery
	err := c.BindJSON(&query)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}
	PlatformID, _ := c.Get("platform_id")

	if err != nil {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}

	err = DB.Model(models.Orders{}).Where(`
		id in ? AND 
		status = 21 AND 
		logistics_status = 0 AND 
		platform_id = ?
	`, query.IDs, PlatformID.(int)).Update("logistics_status", 1).Error

	if err != nil {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}

	g.Response(http.StatusOK, e.Success, nil)
}

func OrderMakeShipmentNo(c *gin.Context) {
	g := Gin{c}
	var query models.OrderListReq
	err := c.BindJSON(&query)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}
	PlatformID, _ := c.Get("platform_id")
	query.PlatformID = PlatformID.(int)
	// 一定要已付款 而且還沒產生寄件編號的
	query.Status = 21
	query.LogisticsStatus = 1
	orders, _ := query.FetchAll()

	if err != nil {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}

	for _, order := range orders {
		order.GetProducts()
		err = money.CreateLogisticsOrder(order)

		if err != nil {
			g.Response(http.StatusBadRequest, e.InvalidParams, err)
			return
		}
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
	query.Status = 21
	query.LogisticsStatus = 21
	orders, err := query.FetchAll()

	if err != nil {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}

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
		params["CVSPaymentNo"] = strings.Join(shipmentNo, ",")
		params["CVSValidationNo"] = ""
		url = "https://logistics.ecpay.com.tw/Express/PrintFAMIC2COrderInfo"
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
