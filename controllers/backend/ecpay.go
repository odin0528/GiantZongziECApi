package backend

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	. "eCommerce/internal/database"
	models "eCommerce/models/backend"

	"github.com/Laysi/go-ecpay-sdk"
	"github.com/gin-gonic/gin"
	"github.com/liudng/godump"
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

	senderMac := params["CheckMacValue"]
	delete(params, "CheckMacValue")
	client := ecpay.NewStageClient(
		ecpay.WithReturnURL(fmt.Sprintf("%s%s", os.Getenv("API_URL"), os.Getenv("ECPAY_PAYMENT_FINISH_URL"))),
	)
	mac := client.GenerateCheckMacValue(params)
	if mac != senderMac {
		c.String(http.StatusBadRequest, "0|Error")
		c.Abort()
	}

	/* if params["SimulatePaid"] == "1" {
		godump.Dump(params)
	} */

	info := QueryTradeInfo(params["MerchantTradeNo"])
	godump.Dump(info.Get("TradeStatus"))

	if info.Get("TradeStatus") == "1" {
		DB.Debug().Model(&models.Orders{}).Where("id = ? and status = 11", strings.Replace(info.Get("MerchantTradeNo"), "GZEC", "", 1)).Update("status", 21)
	}

	c.String(http.StatusOK, "1|OK")
}

func EcpayLogisticsCreate(c *gin.Context) {
	ecpayValue := map[string]string{}
	ecpayValue["GoodsAmount"] = "978"
	ecpayValue["IsCollection"] = "Y"
	ecpayValue["LogisticsSubType"] = "UNIMART"
	ecpayValue["LogisticsType"] = "CVS"
	ecpayValue["MerchantID"] = "2000132"
	ecpayValue["MerchantTradeNo"] = ""
	ecpayValue["MerchantTradeDate"] = time.Now().Format("2006/01/02 15:04:05")
	ecpayValue["ReceiverCellPhone"] = "0958259061"
	ecpayValue["ReceiverName"] = "李晧瑋"
	ecpayValue["ReceiverStoreID"] = "210960"
	ecpayValue["SenderName"] = "李晧瑋"
	ecpayValue["ServerReplyURL"] = fmt.Sprintf("%s/api/backend/ecpay/logistics", os.Getenv("API_URL"))

	encodedParams := fmt.Sprintf(
		"HashKey=%s&%s&HashIV=%s",
		"5294y06JbISpM5x9",
		ecpay.NewECPayValuesFromMap(ecpayValue).Encode(),
		"v77hoKGq4kWxNNIS",
	)

	encodedParams = FormUrlEncode(encodedParams)
	encodedParams = strings.ToLower(encodedParams)
	sum := md5.Sum([]byte(encodedParams))
	checkMac := strings.ToUpper(fmt.Sprintf("%x", sum))

	resp, err := http.Post(
		"https://logistics-stage.ecpay.com.tw/Express/Create",
		"application/x-www-form-urlencoded",
		strings.NewReader(ecpay.NewECPayValuesFromMap(ecpayValue).Encode()+"&CheckMacValue="+checkMac),
	)

	if err != nil {
		log.Fatal(err)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
	response, _ := url.ParseQuery(bodyString)
	godump.Dump(response)

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
