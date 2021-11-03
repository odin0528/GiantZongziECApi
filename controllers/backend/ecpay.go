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
	ecpayBase "github.com/Laysi/go-ecpay-sdk/base"
	ecpayGin "github.com/Laysi/go-ecpay-sdk/gin"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/liudng/godump"
)

func EcpayPaymentFinish(c *gin.Context) {
	data := ecpayBase.OrderResult{}
	err := ecpayGin.ResponseBodyDateTimePatchHelper(c)
	if err != nil {
		fmt.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	if err = c.MustBindWith(&data, binding.FormPost); err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("================")
	godump.Dump(data)

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

	if params["SimulatePaid"] == "1" {
		fmt.Println("================")
		godump.Dump(params)
	}

	info, _, _ := client.QueryTradeInfo(params["MerchantTradeNo"], time.Now())
	fmt.Println("================")
	godump.Dump(info)

	if info.TradeStatus == "1" {
		DB.Model(&models.Orders{}).Where("id = ? and status = 11", strings.Replace(info.MerchantTradeNo, "GZEC", "", 1)).Update("status", 21)
	}

	fmt.Println("1|ok")
	c.String(http.StatusBadRequest, "1|ok")
}
