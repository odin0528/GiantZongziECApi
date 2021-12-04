package backend

import (
	"bytes"
	"crypto/sha256"
	"eCommerce/pkg/e"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	. "eCommerce/internal/database"
	"eCommerce/internal/line"
	models "eCommerce/models/backend"

	"github.com/Laysi/go-ecpay-sdk"
	"github.com/gin-gonic/gin"
	"github.com/liudng/godump"
	log "github.com/sirupsen/logrus"
)

func EcpayPaymentFinish(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	err = c.Request.ParseForm()
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	params := ecpay.ECPayValues{c.Request.PostForm}.ToMap()
	c.Request.Form = nil
	c.Request.PostForm = nil

	log.Println(params)

	senderMac := params["CheckMacValue"]
	delete(params, "CheckMacValue")

	var client *ecpay.Client
	if os.Getenv("ENV") != "production" {
		client = ecpay.NewStageClient(
			ecpay.WithReturnURL(fmt.Sprintf("%s%s", os.Getenv("API_URL"), os.Getenv("ECPAY_PAYMENT_FINISH_URL"))),
			ecpay.WithDebug,
		)
	} else {
		client = ecpay.NewClient(
			os.Getenv("ECPAY_MERCHANT_ID"),
			os.Getenv("ECPAY_MERCHANT_HASH_KEY"),
			os.Getenv("ECPAY_MERCHANT_HASH_IV"),
			fmt.Sprintf("%s%s", os.Getenv("API_URL"), os.Getenv("ECPAY_PAYMENT_FINISH_URL")),
			ecpay.WithDebug,
		)
	}

	mac := client.GenerateCheckMacValue(params)
	if mac != senderMac {
		c.String(http.StatusBadRequest, "0|Error")
		c.Abort()
	}

	/* if params["SimulatePaid"] == "1" {
		godump.Dump(params)
	} */

	info := QueryTradeInfo(params["MerchantTradeNo"])
	godump.Dump(info)

	if info.Get("TradeStatus") == "1" {

		ID, _ := strconv.ParseInt(strings.Replace(info.Get("MerchantTradeNo"), os.Getenv("ECPAY_MERCHANT_TRADE_NO_PREFIX"), "", 1), 10, 0)
		PaymentTypeChargeFee, _ := strconv.ParseFloat(info.Get("PaymentTypeChargeFee"), 32)

		query := models.OrderQuery{
			ID: int(ID),
		}
		order, err := query.FetchForLogistics()
		if err != nil || order.Status != 11 {
			c.String(http.StatusBadRequest, "0|Error")
			c.Abort()
			return
		}

		order.Status = 21
		order.PaymentChargeFee = PaymentTypeChargeFee
		DB.Select("status", "payment_charge_fee").Updates(&order)

		order.GetProducts()

		line.SendOrderNotifyByOrder(order)
	}

	c.String(http.StatusOK, "1|OK")
}

func EcpayLogisticsNotify(c *gin.Context) {
	g := Gin{c}
	body, err := c.GetRawData()
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	err = c.Request.ParseForm()
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	params := ecpay.ECPayValues{c.Request.PostForm}.ToMap()
	c.Request.Form = nil
	c.Request.PostForm = nil

	godump.Dump(params)
	log.Println(params)

	if params["MerchantID"] == os.Getenv("ECPAY_MERCHANT_ID") {
		if tradeNo, ok := params["MerchantTradeNo"]; ok {
			tradeNo = strings.ReplaceAll(tradeNo, os.Getenv("ECPAY_MERCHANT_TRADE_NO_PREFIX"), "")
			id, _ := strconv.ParseInt(tradeNo, 10, 0)

			query := models.OrderQuery{
				ID: int(id),
			}

			order, err := query.FetchForLogistics()

			if err != nil {
				g.Response(http.StatusBadRequest, e.StatusNotFound, err)
				return
			}

			ChangeLogisticsStatus(&order, params)

			err = DB.Select("logistics_msg", "logistics_status", "shipment_no", "logistics_id", "status").Updates(&order).Error
			c.String(http.StatusOK, "1|OK")
		} else {
			c.String(http.StatusBadRequest, "0|Error")
		}
	}
}

func EcpayPaymentTest(c *gin.Context) {
	c.String(http.StatusOK, "1|ok")
}

func QueryTradeInfo(merchantTradeNo string) url.Values {
	timestamp := time.Now().Unix()
	encodedParams := fmt.Sprintf(
		"HashKey=%s&%s&HashIV=%s",
		os.Getenv("ECPAY_MERCHANT_HASH_KEY"),
		fmt.Sprintf("MerchantID=%s&MerchantTradeNo=%s&PlatformID=&TimeStamp=%d", os.Getenv("ECPAY_MERCHANT_ID"), merchantTradeNo, timestamp),
		os.Getenv("ECPAY_MERCHANT_HASH_IV"),
	)
	encodedParams = FormUrlEncode(encodedParams)
	encodedParams = strings.ToLower(encodedParams)
	sum := sha256.Sum256([]byte(encodedParams))
	checkMac := strings.ToUpper(hex.EncodeToString(sum[:]))

	data := url.Values{}

	data.Add("MerchantID", os.Getenv("ECPAY_MERCHANT_ID"))
	data.Add("MerchantTradeNo", merchantTradeNo)
	data.Add("TimeStamp", fmt.Sprintf("%d", timestamp))
	data.Add("CheckMacValue", checkMac)
	data.Add("PlatformID", "")
	resp, err := http.PostForm(os.Getenv("ECPAY_QUERY_TRADE_URL"), data)

	if err != nil {
		log.Fatal(err)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	params, _ := url.ParseQuery(bodyString)

	return params
}

func FormUrlEncode(s string) string {
	s = url.QueryEscape(s)
	//s = strings.ReplaceAll(s, "%2d", "-")
	//s = strings.ReplaceAll(s, "%5f", "_")
	//s = strings.ReplaceAll(s, "%2e", ".")
	s = strings.ReplaceAll(s, "%21", "!")
	s = strings.ReplaceAll(s, "%2A", "*")
	s = strings.ReplaceAll(s, "%28", "(")
	s = strings.ReplaceAll(s, "%29", ")")
	return s
}

func ChangeLogisticsStatus(order *models.Orders, params map[string]string) {
	order.LogisticsID = params["AllPayLogisticsID"]
	order.LogisticsMsg = params["RtnMsg"]
	switch order.Method {
	case 1:
		order.ShipmentNo = params["BookingNote"]
		switch params["RtnCode"] {
		case "300":
			fallthrough
		case "310":
			order.Status = 23
		case "3006":
			order.Status = 51
			order.LogisticsStatus = 110
		case "3001": //轉運中
			order.Status = 51
			order.LogisticsStatus = 110
		case "3003": //配完
			order.Status = 91
			order.LogisticsStatus = 199
		case "5004": //一般單退回
			order.Status = 61
			order.LogisticsStatus = 210
		case "5008": //退貨配完
			order.Status = 92
			order.LogisticsStatus = 299
		default:
			order.LogisticsStatus = 999
		}
	case 2:
		order.ShipmentNo = params["CVSPaymentNo"] + params["CVSValidationNo"]
		switch params["RtnCode"] {
		case "300":
			fallthrough
		case "310":
			order.Status = 23
		case "2068": //交貨便收件(A門市收到件寄件商品)
			order.Status = 51
			order.LogisticsStatus = 110
		case "2073": //商品配達買家取貨門市
			order.LogisticsStatus = 120
		case "2067":
			order.Status = 91
			order.LogisticsStatus = 199
		case "2074":
			order.Status = 61
			order.LogisticsStatus = 210
		case "2072":
			order.Status = 61
			order.LogisticsStatus = 220
		case "2069":
			order.Status = 61
			order.LogisticsStatus = 220
		case "2077":
			order.Status = 92
			order.LogisticsStatus = 299
		default:
			order.LogisticsStatus = 999
		}
	case 3:
		order.ShipmentNo = params["CVSPaymentNo"]
		switch params["RtnCode"] {
		case "300":
			fallthrough
		case "310":
			order.Status = 23
		case "3024": //貨件已至物流中心
			fallthrough
		case "3032": //賣家已到門市寄件
			order.Status = 51
			order.LogisticsStatus = 110
		case "3018": //到店尚未取貨，簡訊通知取件
			order.Status = 51
			order.LogisticsStatus = 120
		case "3022":
			order.Status = 91
			order.LogisticsStatus = 199
		case "3025":
			order.Status = 61
			order.LogisticsStatus = 210
		case "3020":
			order.Status = 61
			order.LogisticsStatus = 210
		case "3019":
			order.Status = 61
			order.LogisticsStatus = 220
		case "3023":
			order.Status = 92
			order.LogisticsStatus = 299
		default:
			order.LogisticsStatus = 999
		}
	case 4:
		order.ShipmentNo = params["CVSPaymentNo"]
		switch params["RtnCode"] {
		case "300":
			fallthrough
		case "310":
			fallthrough
		case "2101":
			order.Status = 23
		case "2041": //物流中心理貨中
			fallthrough
		case "3032": //賣家已到門市寄件
			fallthrough
		case "2068":
			order.Status = 51
			order.LogisticsStatus = 110
		case "2063": //門市配達
			order.Status = 51
			order.LogisticsStatus = 120
		case "2067":
			order.Status = 91
			order.LogisticsStatus = 199
		case "2074":
			order.Status = 61
			order.LogisticsStatus = 210
		case "2069":
			order.Status = 61
			order.LogisticsStatus = 220
		case "2070":
			order.Status = 92
			order.LogisticsStatus = 299
		default:
			order.LogisticsStatus = 999
		}
	case 5:
		order.ShipmentNo = params["CVSPaymentNo"]
		switch params["RtnCode"] {
		case "300":
			order.Status = 23
		case "2030": // 物流中心驗收成功
			fallthrough
		case "3032":
			order.Status = 51
			order.LogisticsStatus = 110
		case "2073": //商品配達買家取貨門市
			order.Status = 51
			order.LogisticsStatus = 120
		case "3022":
			order.Status = 91
			order.LogisticsStatus = 199
		case "2074":
			order.Status = 61
			order.LogisticsStatus = 210
		case "2072":
			order.Status = 61
			order.LogisticsStatus = 220
		case "3023":
			order.Status = 92
			order.LogisticsStatus = 299
		default:
			order.LogisticsStatus = 999
		}
	}
}
