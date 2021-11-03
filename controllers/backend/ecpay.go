package backend

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	. "eCommerce/internal/database"
	models "eCommerce/models/backend"

	"github.com/Laysi/go-ecpay-sdk"
	ecpayGin "github.com/Laysi/go-ecpay-sdk/gin"

	"github.com/gin-gonic/gin"
	"github.com/liudng/godump"
)

func EcpayPaymentFinish(c *gin.Context) {
	ecpayGin.ResponseBodyDateTimePatchHelper(c)

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

	if params["SimulatePaid"] == "1" {
		godump.Dump(params)
	}

	info, _, _ := client.QueryTradeInfo(params["MerchantTradeNo"], time.Now())
	godump.Dump(info)

	if info.TradeStatus == "1" {
		DB.Model(&models.Orders{}).Where("id = ? and status = 11", strings.Replace(info.MerchantTradeNo, "GZEC", "", 1)).Update("status", 21)
	}

	fmt.Println("1|ok")
	c.String(http.StatusBadRequest, "1|ok")
}
