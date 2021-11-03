package backend

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	. "eCommerce/internal/database"
	models "eCommerce/models/backend"

	"github.com/Laysi/go-ecpay-sdk"
	"github.com/gin-gonic/gin"
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
		ecpay.WithDebug,
	)
	mac := client.GenerateCheckMacValue(params)
	if mac != senderMac {
		c.String(http.StatusBadRequest, "0|Error")
		c.Abort()
	}

	info, _, _ := client.QueryTradeInfo(params["MerchantTradeNo"], time.Now())

	if info.TradeStatus == "1" {
		DB.Debug().Model(&models.Orders{}).Where("id = ? and status = 11 and deleted_at = 0", strings.Replace(info.MerchantTradeNo, "GZEC", "", 1)).Update("status", 21)
	}

	fmt.Println("1|ok")
}
