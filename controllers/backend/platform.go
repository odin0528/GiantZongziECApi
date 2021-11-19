package backend

import (
	"eCommerce/pkg/e"
	"fmt"
	"net/http"
	"strings"

	. "eCommerce/internal/database"
	"eCommerce/internal/uploader"
	models "eCommerce/models/backend"

	"github.com/gin-gonic/gin"
)

func PlatformFetch(c *gin.Context) {
	g := Gin{c}
	PlatformID, _ := c.Get("platform_id")
	platform := &models.Platform{
		ID: PlatformID.(int),
	}

	platform.Fetch()

	g.Response(http.StatusOK, e.Success, platform)
}

func PlatformUpdate(c *gin.Context) {
	g := Gin{c}
	var platform *models.Platform
	PlatformID, _ := c.Get("platform_id")
	c.BindJSON(&platform)
	platform.ID = PlatformID.(int)

	if strings.Index(platform.LogoUrl, ",") > 0 {
		filename := fmt.Sprintf("/upload/%08d/logo", PlatformID.(int))
		platform.LogoUrl = uploader.Thumbnail(filename, platform.LogoUrl, 2048)
	}

	if strings.Index(platform.IconUrl, ",") > 0 {
		filename := fmt.Sprintf("/upload/%08d/favicon", PlatformID.(int))
		platform.IconUrl = uploader.Thumbnail(filename, platform.IconUrl, 2048)
	}

	DB.Select("title", "description", "logo_url", "icon_url", "fb_messenger_enabled", "fb_pixel").Updates(&platform)

	g.Response(http.StatusOK, e.Success, platform)
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

func PlatformPaymentUpdate(c *gin.Context) {
	g := Gin{c}
	var payment *models.PlatformPayment
	PlatformID, _ := c.Get("platform_id")
	c.BindJSON(&payment)

	DB.Where("platform_id = ?", PlatformID.(int)).Updates(&payment)

	g.Response(http.StatusOK, e.Success, nil)
}

func PlatformLogisticsFetch(c *gin.Context) {
	g := Gin{c}
	PlatformID, _ := c.Get("platform_id")
	logistics := &models.PlatformLogistics{
		PlatformID: PlatformID.(int),
	}

	logistics.Fetch()

	g.Response(http.StatusOK, e.Success, logistics)
}

func PlatformLogisticsUpdate(c *gin.Context) {
	g := Gin{c}
	var logistics *models.PlatformLogistics
	PlatformID, _ := c.Get("platform_id")
	c.BindJSON(&logistics)

	DB.Debug().Where("platform_id = ?", PlatformID.(int)).Updates(&logistics)

	g.Response(http.StatusOK, e.Success, nil)
}
