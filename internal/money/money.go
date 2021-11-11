package money

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
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

type QueryLogisticsInfoReq struct {
	PlatformID string
	MerchantID string
	RqHeader   EcpayReqHeader
	Data       string
}

type EcpayReqHeader struct {
	Timestamp int64
	Revision  string
}

type EcpayReqBody struct {
	MerchantID      string
	LogisticsID     string
	MerchantTradeNo string
}

type QueryLogisticsInfoRes struct {
	PlatformID string
	MerchantID string
	RqHeader   EcpayReqHeader
	Data       string
	TransCode  int
	TransMsg   string
}

func CreateLogisticsOrder(order models.Orders) error {
	/* goodsName := []string{}
	for _, product := range order.Products {
		goodsName = append(goodsName, FilterGoodsName(product.Title))
	} */

	ecpayValue := map[string]string{}
	ecpayValue["MerchantID"] = os.Getenv("ECPAY_MERCHANT_ID")
	ecpayValue["GoodsAmount"] = fmt.Sprintf("%d", int(order.Total))
	ecpayValue["MerchantTradeDate"] = time.Now().Format("2006/01/02 15:04:05")
	ecpayValue["MerchantTradeNo"] = fmt.Sprintf("%s%d", os.Getenv("ECPAY_MERCHANT_TRADE_NO_PREFIX"), order.ID)
	ecpayValue["ReceiverCellPhone"] = order.Phone
	ecpayValue["ReceiverName"] = order.Fullname
	ecpayValue["SenderName"] = "李晧瑋"
	ecpayValue["SenderCellPhone"] = "0958259061"
	ecpayValue["ServerReplyURL"] = fmt.Sprintf("%s/api/backend/ecpay/logistics", os.Getenv("API_URL"))

	goodsName := []rune(FilterGoodsName(order.Products[0].Title))
	if len(goodsName) > 25 {
		ecpayValue["GoodsName"] = string(goodsName[:25])
	} else {
		ecpayValue["GoodsName"] = string(goodsName)
	}

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
		ecpayValue["ReceiverZipCode"] = order.ZipCode
		ecpayValue["ReceiverAddress"] = order.County + order.District + order.Address
		ecpayValue["Temperature"] = "0001"
		ecpayValue["Distance"] = "00"
		ecpayValue["Specification"] = "0004"
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
		fmt.Sprintf("%s%s", os.Getenv("ECPAY_LOGISTICS_URL"), "/Express/Create"),
		"application/x-www-form-urlencoded",
		strings.NewReader(ecpay.NewECPayValuesFromMap(ecpayValue).Encode()+"&CheckMacValue="+checkMac),
	)

	if err != nil {
		log.Fatal(err)
		return err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	bodyString := strings.Split(string(bodyBytes), "|")

	if bodyString[0] != "1" {
		return errors.New(fmt.Sprintf("建立物流訂單失敗 回傳訊息為：%s", bodyString[1]))
	}

	return nil
}

func QueryLogisticsInfo(allPayLogisticsID string) (info url.Values, err error) {
	ecpayValue := map[string]string{}
	ecpayValue["MerchantID"] = os.Getenv("ECPAY_MERCHANT_ID")
	ecpayValue["AllPayLogisticsID"] = allPayLogisticsID
	ecpayValue["TimeStamp"] = fmt.Sprintf("%d", time.Now().Unix())
	ecpayValue["PlatformID"] = ""

	checkMac := MakeLogisticsCheckMac(ecpayValue)

	resp, err := http.Post(
		fmt.Sprintf("%s%s", os.Getenv("ECPAY_LOGISTICS_URL"), "/Helper/QueryLogisticsTradeInfo/V2"),
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

func QueryLogisticsInfoV2(allPayLogisticsID string) (info url.Values, err error) {
	ecpayValue := QueryLogisticsInfoReq{}
	ecpayValue.PlatformID = ""
	ecpayValue.MerchantID = os.Getenv("ECPAY_MERCHANT_ID")
	ecpayValue.RqHeader.Timestamp = time.Now().Unix()
	ecpayValue.RqHeader.Revision = "1.0.0"

	ecpayBody := EcpayReqBody{}
	ecpayBody.MerchantID = os.Getenv("ECPAY_MERCHANT_ID")
	ecpayBody.LogisticsID = allPayLogisticsID

	bodyData, _ := json.Marshal(ecpayBody)

	key := []byte(os.Getenv("ECPAY_MERCHANT_HASH_KEY"))
	iv := []byte(os.Getenv("ECPAY_MERCHANT_HASH_IV"))
	plaintext := []byte(ecpay.FormUrlEncode(string(bodyData)))

	plaintext = PKCS7Padding(plaintext)
	ciphertext := make([]byte, len(plaintext))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	ecpayValue.Data = base64.StdEncoding.EncodeToString(ciphertext)
	postData, _ := json.Marshal(ecpayValue)

	resp, err := http.Post(
		fmt.Sprintf("%s%s", os.Getenv("ECPAY_LOGISTICS_URL"), "/Express/v2/QueryLogisticsTradeInfo"),
		"application/json",
		bytes.NewBuffer(postData),
	)

	if err != nil {
		log.Fatal(err)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	var res QueryLogisticsInfoRes
	json.Unmarshal(bodyBytes, &res)

	decode, _ := base64.StdEncoding.DecodeString(res.Data)
	println(decode)
	godump.Dump(res)

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

	// println(encodedParams)

	encodedParams = ecpay.FormUrlEncode(encodedParams)
	encodedParams = strings.ToLower(encodedParams)

	return encodedParams
}

func MakeLogisticsCheckMac(params map[string]string) string {
	encodedParams := MakeQueryString(params)
	// println(encodedParams)
	sum := md5.Sum([]byte(encodedParams))
	checkMac := strings.ToUpper(fmt.Sprintf("%x", sum))

	return checkMac
}

func GeneratePrintShipmentCheckMac(ids []string) string {
	ecpayValue := map[string]string{}
	ecpayValue["MerchantID"] = os.Getenv("ECPAY_MERCHANT_ID")
	ecpayValue["AllPayLogisticsID"] = strings.Join(ids, ",")
	ecpayValue["MerchantTradeNo"] = ""
	ecpayValue["IsPreView"] = "True"

	checkMac := MakeLogisticsCheckMac(ecpayValue)

	return checkMac
}

func FilterGoodsName(s string) string {
	s = strings.ReplaceAll(s, "^", "＾")
	s = strings.ReplaceAll(s, "'", "＇")
	s = strings.ReplaceAll(s, "`", "｀")
	s = strings.ReplaceAll(s, "!", "！")
	s = strings.ReplaceAll(s, "@", "＠")
	s = strings.ReplaceAll(s, "#", "＃")
	s = strings.ReplaceAll(s, "%", "％")
	s = strings.ReplaceAll(s, "&", "＆")
	s = strings.ReplaceAll(s, "*", "＊")
	s = strings.ReplaceAll(s, "+", "＋")
	s = strings.ReplaceAll(s, "\\", "／")
	s = strings.ReplaceAll(s, "\"", "＂")
	s = strings.ReplaceAll(s, "<", "＜")
	s = strings.ReplaceAll(s, ">", "＞")
	s = strings.ReplaceAll(s, "|", "｜")
	s = strings.ReplaceAll(s, "_", "＿")
	s = strings.ReplaceAll(s, "[", "［")
	s = strings.ReplaceAll(s, "]", "］")
	s = strings.ReplaceAll(s, "(", "（")
	s = strings.ReplaceAll(s, ")", "）")
	s = strings.ReplaceAll(s, " ", "　")
	s = strings.ReplaceAll(s, "-", "－")

	return s
}

func PKCS7Padding(ciphertext []byte) []byte {
	padding := aes.BlockSize - len(ciphertext)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
