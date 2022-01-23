package backend

import (
	"eCommerce/pkg/e"
	"net/http"

	models "eCommerce/models/backend"

	"github.com/gin-gonic/gin"
)

func BrandFetchAll(c *gin.Context) {
	g := Gin{c}
	var query models.BrandQuery
	PlatformID, _ := c.Get("platform_id")
	query.PlatformID = PlatformID.(int)
	brands, err := query.FetchAll()

	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	g.Response(http.StatusOK, e.Success, brands)
}

func BrandFetchByPage(c *gin.Context) {
	g := Gin{c}
	var query models.BrandQuery
	g.C.BindJSON(&query)
	PlatformID, _ := c.Get("platform_id")
	query.PlatformID = PlatformID.(int)
	brands, pagination := query.Fetch()

	g.PaginationResponse(http.StatusOK, e.Success, brands, pagination)
}
