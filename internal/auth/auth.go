package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
	// "ec/internal/redis"
)

func AuthRequred(c *gin.Context) {
	fmt.Println(c.Request.Header["Authorization"])

	c.Set("customer_id", 1)
	// rdb.Get(ctx, "odin")

	// redis.Set("odin", "cool", 30)

	// text, _ := redis.Get("odin")
	// fmt.Println(text)
}
