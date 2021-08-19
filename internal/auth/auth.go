package auth

import (
	"regexp"

	"github.com/gin-gonic/gin"

	// "ec/internal/redis"

	"eCommerce/models/backend"
	"eCommerce/models/frontend"
	"eCommerce/pkg/e"
)

func AuthRequred(c *gin.Context) {

	if c.Request.Header.Get("Authorization") == "" {
		c.Abort()
		c.JSON(200, gin.H{
			"http_status": 401,
			"code":        401,
			"msg":         e.GetMsg(401),
			"data":        nil,
		})
		return
	}

	token := backend.AdminToken{
		Token: c.Request.Header.Get("Authorization"),
	}

	token.Fetch()

	if token.AdminID == 0 {
		c.Abort()
		c.JSON(200, gin.H{
			"http_status": 401,
			"code":        401,
			"msg":         e.GetMsg(401),
			"data":        nil,
		})
		return
	}

	query := backend.AdminQuery{
		ID: token.AdminID,
	}

	admin := query.Fetch()

	c.Set("customer_id", admin.CustomerID)
}

func GetCustomerID(c *gin.Context) {
	r, _ := regexp.Compile("^([a-zA-Z0-9\\.]*).*$")
	match := r.FindAllStringSubmatch(c.Request.Header["Hostname"][0], 1)

	query := &frontend.CustomerQuery{}
	query.Hostname = match[0][1]
	customer := query.Fetch()

	if customer.ID == 0 {
		c.Abort()
		c.JSON(404, gin.H{
			"http_status": 404,
			"code":        404,
			"msg":         e.GetMsg(404),
			"data":        nil,
		})
	}

	c.Set("customer_id", customer.ID)
}
