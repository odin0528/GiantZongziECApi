package backend

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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
		ecpay.WithReturnURL(os.Getenv("API_URL")),
		ecpay.WithOrderResultURL(fmt.Sprintf(os.Getenv("ECPAY_CLIENT_RETURN_URL"), c.Request.Header["Hostname"][0], "%2Fcheckout%2Ffinish")),
		ecpay.WithDebug,
	)
	mac := client.GenerateCheckMacValue(params)
	fmt.Println(mac)
	if mac != senderMac {
		c.String(http.StatusBadRequest, "0|Error")
		c.Abort()
	}
	fmt.Println("1|ok")
}
