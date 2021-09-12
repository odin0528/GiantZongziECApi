package frontend

import (
	"eCommerce/pkg/e"
	"net/http"

	"github.com/gin-gonic/gin"

	models "eCommerce/models/frontend"
)

func PlatformFetch(c *gin.Context) {
	g := Gin{c}
	platform, _ := c.Get("platform")
	Platform := platform.(models.Platform)
	menu := Platform.GetMenu()
	promotions := Platform.GetPromotions()
	g.Response(http.StatusOK, e.Success, map[string]interface{}{"info": platform, "menu": menu, "promotions": promotions})
}
