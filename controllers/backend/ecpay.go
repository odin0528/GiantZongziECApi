package backend

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func EcpayPaymentFinish(c *gin.Context) {
	checkMacValue := c.PostForm("CheckMacValue")
	fmt.Println(checkMacValue)
}
