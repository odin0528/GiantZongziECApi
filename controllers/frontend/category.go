package frontend

import (
	"eCommerce/pkg/e"
	"net/http"

	models "eCommerce/models/frontend"

	"github.com/gin-gonic/gin"
)

func CategoryList(c *gin.Context) {
	g := Gin{c}
	var query models.CategoryQuery
	err := c.ShouldBindUri(&query)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}
	PlatformID, _ := c.Get("platform_id")
	query.PlatformID = PlatformID.(int)
	categories := query.FetchAll()

	var breadcrumbs []models.Category
	query.GetBreadcrumbs(&breadcrumbs)

	g.Response(http.StatusOK, e.Success, map[string]interface{}{"categories": categories, "breadcrumbs": breadcrumbs})
}

func CategoryChildList(c *gin.Context) {
	g := Gin{c}
	var query models.CategoryQuery
	err := c.ShouldBindUri(&query)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}
	PlatformID, _ := c.Get("platform_id")
	query.PlatformID = PlatformID.(int)
	categories := query.FetchAll()

	g.Response(http.StatusOK, e.Success, categories)
}
