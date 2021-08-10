package auth

import (
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"

	// "ec/internal/redis"

	"eCommerce/models/frontend"
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

	c.Set("customer_id", customer.ID)
}
