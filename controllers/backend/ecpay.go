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
		DB.Debug().Model(&models.Orders{}).Where("id = ? and status = 11", strings.Replace(info.Get("MerchantTradeNo"), os.Getenv("ECPAY_MERCHANT_TRADE_NO_PREFIX"), "", 1)).Update("status", 21)
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
			godump.Dump(order)

			err = DB.Select("logistics_msg", "logistics_status", "shipment_no", "logistics_id").Updates(&order).Error
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
			order.LogisticsStatus = 2
		case "3006":
			order.Status = 31
			order.LogisticsStatus = 3
		case "3003":
			order.Status = 91
			order.LogisticsStatus = 5
		case "5004":
			order.LogisticsStatus = 6
		case "5008":
			order.Status = 92
			order.LogisticsStatus = 8
		default:
			order.LogisticsStatus = 99
		}
	case 2:
		order.ShipmentNo = params["CVSPaymentNo"] + params["CVSValidationNo"]
		switch params["RtnCode"] {
		case "300":
			fallthrough
		case "310":
			order.LogisticsStatus = 2
		case "2068":
			order.Status = 31
			order.LogisticsStatus = 3
		case "2073":
			order.LogisticsStatus = 4
		case "2067":
			order.Status = 91
			order.LogisticsStatus = 5
		case "2074":
			order.LogisticsStatus = 6
		case "2069":
			order.LogisticsStatus = 7
		case "2077":
			order.Status = 92
			order.LogisticsStatus = 8
		default:
			order.LogisticsStatus = 99
		}
	case 3:
		order.ShipmentNo = params["CVSPaymentNo"]
		switch params["RtnCode"] {
		case "300":
			fallthrough
		case "310":
			order.LogisticsStatus = 2
		case "3032":
			order.Status = 31
			order.LogisticsStatus = 3
		case "3018":
			order.LogisticsStatus = 4
		case "3022":
			order.Status = 91
			order.LogisticsStatus = 5
		case "3020":
			order.LogisticsStatus = 6
		case "3019":
			order.LogisticsStatus = 7
		case "3023":
			order.Status = 92
			order.LogisticsStatus = 8
		default:
			order.LogisticsStatus = 99
		}
	case 4:
		order.ShipmentNo = params["CVSPaymentNo"]
		switch params["RtnCode"] {
		case "300":
			fallthrough
		case "310":
			order.LogisticsStatus = 2
		case "2068":
			order.Status = 31
			order.LogisticsStatus = 3
		case "2073":
			order.LogisticsStatus = 4
		case "2067":
			order.Status = 91
			order.LogisticsStatus = 5
		case "2074":
			order.LogisticsStatus = 6
		case "2069":
			order.LogisticsStatus = 7
		case "2070":
			order.Status = 92
			order.LogisticsStatus = 8
		default:
			order.LogisticsStatus = 99
		}
	case 5:
		order.ShipmentNo = params["CVSPaymentNo"]
		switch params["RtnCode"] {
		case "300":
			order.LogisticsStatus = 2
		case "2030":
			order.Status = 31
			order.LogisticsStatus = 3
		case "2073":
			order.LogisticsStatus = 4
		case "3022":
			order.Status = 91
			order.LogisticsStatus = 5
		case "2074":
			order.LogisticsStatus = 6
		case "2072":
			order.LogisticsStatus = 7
		case "3023":
			order.Status = 92
			order.LogisticsStatus = 8
		default:
			order.LogisticsStatus = 99
		}
	}
}
