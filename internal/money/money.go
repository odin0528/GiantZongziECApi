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

func CreateLogisticsOrder(order models.Orders) (info url.Values, err error) {
	goodsName := []string{}
	for _, product := range order.Products {
		goodsName = append(goodsName, product.Title)
	}

	ecpayValue := map[string]string{}
	ecpayValue["MerchantID"] = os.Getenv("ECPAY_MERCHANT_ID")
	ecpayValue["GoodsAmount"] = fmt.Sprintf("%d", int(order.Total))
	ecpayValue["MerchantTradeDate"] = time.Now().Format("2006/01/02 15:04:05")
	ecpayValue["MerchantTradeNo"] = fmt.Sprintf("%s%d%d", os.Getenv("ECPAY_MERCHANT_TRADE_NO_PREFIX"), order.ID, time.Now().Unix())
	ecpayValue["ReceiverCellPhone"] = order.Phone
	ecpayValue["ReceiverName"] = order.Fullname
	ecpayValue["SenderName"] = "李晧瑋"
	ecpayValue["SenderCellPhone"] = "0958259061"
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
		ecpayValue["LogisticsSubType"] = "TCAT"
		ecpayValue["SenderZipCode"] = "235"
		ecpayValue["SenderAddress"] = "新北市中和區中正路753號7樓"
		ecpayValue["ReceiverZipCode"] = "104"
		ecpayValue["ReceiverAddress"] = order.Address
		ecpayValue["Temperature"] = "0001"
		ecpayValue["Distance"] = "00"
		ecpayValue["Specification"] = "0004"
	case 2:
		ecpayValue["LogisticsType"] = "CVS"
		ecpayValue["LogisticsSubType"] = "UNIMART"
		ecpayValue["ReceiverStoreID"] = order.StoreID
	case 3:
		ecpayValue["LogisticsType"] = "CVS"
		ecpayValue["LogisticsSubType"] = "FAMI"
		ecpayValue["ReceiverStoreID"] = order.StoreID
	case 4:
		ecpayValue["LogisticsType"] = "CVS"
		ecpayValue["LogisticsSubType"] = "HILIFE"
		ecpayValue["ReceiverStoreID"] = order.StoreID
	case 5:
		ecpayValue["LogisticsType"] = "CVS"
		ecpayValue["LogisticsSubType"] = "OKMART"
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
	bodyString := strings.Split(string(bodyBytes), "|")

	if bodyString[0] != "1" {
		err = errors.New(fmt.Sprintf("建立物流訂單失敗 回傳訊息為：%s", bodyString[1]))
		return
	}

	response, _ := url.ParseQuery(bodyString[1])
	godump.Dump(response)
	info, err = QueryLogisticsInfo(response.Get("AllPayLogisticsID"))

	godump.Dump(info)

	if err != nil {
		err = errors.New(fmt.Sprintf("訂單查詢失敗 回傳代碼為：%s", response.Get("RtnCode")))
		return
	}

	return
}

func QueryLogisticsInfo(allPayLogisticsID string) (info url.Values, err error) {
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

	if bodyString[0:1] == "0" {
		err = errors.New(fmt.Sprintf("建立物流訂單失敗 回傳訊息為：%s", bodyString[2:]))
		return
	}

	bodyString = strings.Replace(bodyString, "1|", "", 1)
	info, _ = url.ParseQuery(bodyString)
	err = nil

	return
}

func MakeQueryString(params map[string]string) string {
	encodedParams := fmt.Sprintf(
		"HashKey=%s&%s&HashIV=%s",
		os.Getenv("ECPAY_MERCHANT_HASH_KEY"),
		ecpay.NewECPayValuesFromMap(params).Encode(),
		os.Getenv("ECPAY_MERCHANT_HASH_IV"),
	)

	println(encodedParams)

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
