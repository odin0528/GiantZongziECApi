package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
	// "eCommerce/internal/redis"
)

func AuthRequred(c *gin.Context) {
	fmt.Println(c.Request.Header["Authorization"])
	// rdb.Get(ctx, "odin")

	// redis.Set("odin", "cool", 30)

	// text, _ := redis.Get("odin")
	// fmt.Println(text)
}
