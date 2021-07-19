package backend

import (
	"eCommerce/pkg/e"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPagesList(c *gin.Context) {
	g := Gin{c}
	g.Response(http.StatusOK, e.Success, "YAYAYAYA")
}
