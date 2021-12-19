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
	logistics := Platform.GetLogistics()
	payment := Platform.GetPayments()
	categories := Platform.GetCategory()
	g.Response(http.StatusOK, e.Success, map[string]interface{}{
		"info":       platform,
		"menu":       menu,
		"promotions": promotions,
		"payment":    payment,
		"logistics":  logistics,
		"categories": categories,
	})
}

func PlatformPaymentFetch(c *gin.Context) {
	g := Gin{c}
	PlatformID, _ := c.Get("platform_id")
	payment := &models.PlatformPayment{
		PlatformID: PlatformID.(int),
	}

	payment.Fetch()

	g.Response(http.StatusOK, e.Success, payment)
}
