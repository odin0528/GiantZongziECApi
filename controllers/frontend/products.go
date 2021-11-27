package frontend

import (
	"eCommerce/pkg/e"
	"net/http"

	models "eCommerce/models/frontend"

	"github.com/gin-gonic/gin"
)

func GetProductsByCategoryID(c *gin.Context) {
	g := Gin{c}
	var req models.ProductQuery
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}
	PlatformID, _ := c.Get("platform_id")
	req.PlatformID = PlatformID.(int)
	products, pagination := req.FetchAll()

	for index := range products {
		products[index].GetPhotos()
	}

	g.PaginationResponse(http.StatusOK, e.Success, products, pagination)
}

func ProductFetch(c *gin.Context) {
	g := Gin{c}
	var query models.ProductQuery
	err := c.ShouldBindUri(&query)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}
	PlatformID, _ := c.Get("platform_id")
	query.PlatformID = PlatformID.(int)
	products := query.Fetch()
	products.GetPhotos()
	products.GetStyle()
	products.GetSubStyle()
	products.GetStyleTable()
	related := products.GetRelated()

	for index := range related {
		related[index].GetPhotos()
	}

	g.Response(http.StatusOK, e.Success, map[string]interface{}{"products": products, "related": related})
}
