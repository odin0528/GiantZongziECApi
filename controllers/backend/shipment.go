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

func ShipmentList(c *gin.Context) {
	g := Gin{c}
	var req models.OrderListReq
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}
	PlatformID, _ := c.Get("platform_id")
	AdminID, _ := c.Get("admin_id")
	req.PlatformID = PlatformID.(int)
	req.PickerID = AdminID.(int)
	orders, pagination := req.FetchAll()

	if req.WithoutProducts != true {
		for index := range orders {
			orders[index].GetProducts()
		}
	}

	g.PaginationResponse(http.StatusOK, e.Success, orders, pagination)
}

func ShipmentSend(c *gin.Context) {
	g := Gin{c}
	var query models.BatchOrderQuery
	err := c.BindJSON(&query)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	PlatformID, _ := c.Get("platform_id")
	err = DB.Model(models.Orders{}).Where(`
		id in ? AND 
		status = 23 AND 
		platform_id = ?
	`, query.IDs, PlatformID.(int)).Updates(map[string]interface{}{"status": 24, "logistics_status": 100}).Error

	if err != nil {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}

	g.Response(http.StatusOK, e.Success, nil)
}

func ShipmentMakeNo(c *gin.Context) {
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
	query.Status = []int{22}
	orders, _ := query.FetchAll()

	if err != nil {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}

	for _, order := range orders {
		order.GetProducts()
		err = money.CreateLogisticsOrder(order)

		if err != nil {
			fmt.Println(err)
			g.Response(http.StatusOK, e.StatusInternalServerError, err.Error())
			return
		}

		err = DB.Model(models.Orders{}).
			Where(`id = ? AND status = 22 AND platform_id = ?`, order.ID, PlatformID.(int)).
			Update("status", 23).Error

		if err != nil {
			g.Response(http.StatusBadRequest, e.StatusNotFound, err)
			return
		}
	}

	g.Response(http.StatusOK, e.Success, nil)
}

func ShipmentPrint(c *gin.Context) {
	g := Gin{c}
	var query models.BatchOrderQuery
	err := c.BindJSON(&query)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}
	PlatformID, _ := c.Get("platform_id")
	query.PlatformID = PlatformID.(int)
	query.Status = 23
	orders, err := query.FetchAll()

	if err != nil || len(orders) == 0 {
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
		url = fmt.Sprintf("%s/helper/PrintTradeDocument", os.Getenv("ECPAY_LOGISTICS_URL"))
	case 2:
		params["CVSPaymentNo"] = strings.Join(paymentNo, ",")
		params["CVSValidationNo"] = strings.Join(validationNo, ",")
		url = fmt.Sprintf("%s/Express/PrintUniMartC2COrderInfo", os.Getenv("ECPAY_LOGISTICS_URL"))
	case 3:
		params["CVSPaymentNo"] = strings.Join(shipmentNo, ",")
		params["CVSValidationNo"] = ""
		url = fmt.Sprintf("%s/Express/PrintFAMIC2COrderInfo", os.Getenv("ECPAY_LOGISTICS_URL"))
	case 4:
		params["MerchantTradeNo"] = strings.Join(tradeNo, ",")
		url = fmt.Sprintf("%s/helper/PrintTradeDocument", os.Getenv("ECPAY_LOGISTICS_URL"))
	case 5:
		params["CVSPaymentNo"] = strings.Join(paymentNo, ",")
		params["CVSValidationNo"] = strings.Join(validationNo, ",")
		url = fmt.Sprintf("%s/Express/PrintOKMARTC2COrderInfo", os.Getenv("ECPAY_LOGISTICS_URL"))
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
