package auth

import (
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"

	// "ec/internal/redis"

	"eCommerce/models/frontend"
	"eCommerce/pkg/e"
)

func AuthRequred(c *gin.Context) {
	fmt.Println(c.Request.Header["Authorization"])

	c.Set("customer_id", 1)
	// rdb.Get(ctx, "odin")

	// redis.Set("odin", "cool", 30)

	// text, _ := redis.Get("odin")
	// fmt.Println(text)
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
