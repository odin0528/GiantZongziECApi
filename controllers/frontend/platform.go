package frontend

import (
	"eCommerce/pkg/e"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PlatformFetch(c *gin.Context) {
	g := Gin{c}
	platform, _ := c.Get("platform")
	g.Response(http.StatusOK, e.Success, platform)
}
