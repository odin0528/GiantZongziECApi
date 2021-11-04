package money

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	models "eCommerce/models/backend"

	"github.com/Laysi/go-ecpay-sdk"
	"github.com/liudng/godump"
)

func CreateLogisticsOrder(order models.Orders) (shipmentNo string, err error) {
	goodsName := []string{}
	for _, product := range order.Products {
		goodsName = append(goodsName, product.Title)
	}

	ecpayValue := map[string]string{}
	ecpayValue["MerchantID"] = os.Getenv("ECPAY_MERCHANT_ID")
	ecpayValue["GoodsAmount"] = fmt.Sprintf("%f", order.Total)
	ecpayValue["MerchantTradeDate"] = time.Now().Format("2006/01/02 15:04:05")
	ecpayValue["MerchantTradeNo"] = fmt.Sprintf("GZEC%d", order.ID)
	ecpayValue["ReceiverCellPhone"] = order.Phone
	ecpayValue["ReceiverName"] = order.Fullname
	ecpayValue["SenderName"] = "李晧瑋"
	ecpayValue["ServerReplyURL"] = fmt.Sprintf("%s/api/backend/ecpay/logistics", os.Getenv("API_URL"))
	ecpayValue["GoodsName"] = strings.Join(goodsName, "#")

	if order.Payment == 2 && order.Method != 1 {
		ecpayValue["IsCollection"] = "Y"
	} else {
		ecpayValue["IsCollection"] = "N"
	}

	switch order.Method {
	case 1:
		ecpayValue["LogisticsType"] = "HOME"
	case 2:
		ecpayValue["LogisticsType"] = "CVS"
		ecpayValue["LogisticsSubType"] = "UNIMARTC2C"
		ecpayValue["ReceiverStoreID"] = order.StoreID
	case 3:
		ecpayValue["LogisticsType"] = "CVS"
		ecpayValue["LogisticsSubType"] = "FAMIC2C"
		ecpayValue["ReceiverStoreID"] = order.StoreID
	case 4:
		ecpayValue["LogisticsType"] = "CVS"
		ecpayValue["LogisticsSubType"] = "HILIFEC2C"
		ecpayValue["ReceiverStoreID"] = order.StoreID
	case 5:
		ecpayValue["LogisticsType"] = "CVS"
		ecpayValue["LogisticsSubType"] = "OKMARTC2C"
		ecpayValue["ReceiverStoreID"] = order.StoreID
	}

	checkMac := MakeLogisticsCheckMac(ecpayValue)

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
	bodyString = strings.Replace(bodyString, "1|", "", 1)
	response, _ := url.ParseQuery(bodyString)

	godump.Dump(response)

	if response.Get("RtnCode") != "100" {
		err = errors.New(fmt.Sprintf("建立物流訂單失敗 回傳代碼為：%s", response.Get("RtnCode")))
		return
	}

	info := QueryLogisticsInfo(response.Get("AllPayLogisticsID"))

	godump.Dump(info)

	shipmentNo = info.Get("ShipmentNo")
	err = nil

	/* if response.Get("RtnCode") != "100" {
		err = errors.New(fmt.Sprintf("建立物流訂單失敗 回傳代碼為：%s", response.Get("RtnCode")))
		return
	} */

	return
}

func QueryLogisticsInfo(allPayLogisticsID string) url.Values {
	ecpayValue := map[string]string{}
	ecpayValue["MerchantID"] = os.Getenv("ECPAY_MERCHANT_ID")
	ecpayValue["AllPayLogisticsID"] = allPayLogisticsID
	ecpayValue["TimeStamp"] = fmt.Sprintf("%d", time.Now().Unix())
	ecpayValue["PlatformID"] = ""

	checkMac := MakeLogisticsCheckMac(ecpayValue)

	resp, err := http.Post(
		"https://logistics-stage.ecpay.com.tw/Helper/QueryLogisticsTradeInfo/V2",
		"application/x-www-form-urlencoded",
		strings.NewReader(ecpay.NewECPayValuesFromMap(ecpayValue).Encode()+"&CheckMacValue="+checkMac),
	)

	if err != nil {
		log.Fatal(err)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	bodyString = strings.Replace(bodyString, "1|", "", 1)
	response, _ := url.ParseQuery(bodyString)

	return response
}

func MakeQueryString(params map[string]string) string {
	encodedParams := fmt.Sprintf(
		"HashKey=%s&%s&HashIV=%s",
		"5294y06JbISpM5x9",
		ecpay.NewECPayValuesFromMap(params).Encode(),
		"v77hoKGq4kWxNNIS",
	)

	encodedParams = ecpay.FormUrlEncode(encodedParams)
	encodedParams = strings.ToLower(encodedParams)

	return encodedParams
}

func MakeLogisticsCheckMac(params map[string]string) string {
	encodedParams := MakeQueryString(params)
	sum := md5.Sum([]byte(encodedParams))
	checkMac := strings.ToUpper(fmt.Sprintf("%x", sum))

	return checkMac
}
