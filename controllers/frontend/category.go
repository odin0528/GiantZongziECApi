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

	var breadcrumbs []models.Category
	query.GetBreadcrumbs(&breadcrumbs)

	if len(breadcrumbs) > 0 {
		query.ID = breadcrumbs[len(breadcrumbs)-1].ParentID
	} else {
		query.ID = -1
	}

	categories := query.FetchAll()

	g.Response(http.StatusOK, e.Success, map[string]interface{}{"categories": categories, "breadcrumbs": breadcrumbs})
}
