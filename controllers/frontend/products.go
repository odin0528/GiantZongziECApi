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
	err := c.ShouldBindUri(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}
	CustomerID, _ := c.Get("customer_id")
	req.CustomerID = CustomerID.(int)
	products, pagination := req.FetchAll()

	for index := range products {
		products[index].GetPhotos()
		products[index].GetPriceRange()
	}

	g.PaginationResponse(http.StatusOK, e.Success, products, pagination)
}
