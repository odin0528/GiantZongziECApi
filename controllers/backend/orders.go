package backend

import (
	"eCommerce/pkg/e"
	"net/http"

	. "eCommerce/internal/database"
	models "eCommerce/models/backend"

	"github.com/gin-gonic/gin"
)

func OrderFetch(c *gin.Context) {
	g := Gin{c}
	var query models.ProductQuery
	err := c.ShouldBindUri(&query)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}
	CustomerID, _ := c.Get("customer_id")
	query.CustomerID = CustomerID.(int)
	product := query.Fetch()
	product.GetPhotos()
	product.GetStyle()
	product.GetSubStyle()
	product.GetStyleTable()

	g.Response(http.StatusOK, e.Success, product)
}

func OrderList(c *gin.Context) {
	g := Gin{c}
	var req models.OrderListReq
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}
	CustomerID, _ := c.Get("customer_id")
	req.CustomerID = CustomerID.(int)
	orders, pagination := req.FetchAll()

	for index := range orders {
		orders[index].GetProducts()
	}

	g.PaginationResponse(http.StatusOK, e.Success, orders, pagination)
}

func OrderNextStep(c *gin.Context) {
	g := Gin{c}
	var query models.OrderQuery
	err := c.BindJSON(&query)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}
	CustomerID, _ := c.Get("customer_id")
	query.CustomerID = CustomerID.(int)
	order, err := query.Fetch()

	if err != nil {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}

	switch order.Status {
	case 11:
		order.Status = 21
	case 21:
		order.Status = 31
	case 31:
		order.Status = 99
	}

	err = DB.Debug().Select("status").Updates(&order).Error

	if err != nil {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}

	g.Response(http.StatusOK, e.Success, nil)
}

func OrderUntreated(c *gin.Context) {
	g := Gin{c}
	var query models.OrderQuery
	CustomerID, _ := c.Get("customer_id")
	query.CustomerID = CustomerID.(int)
	count := query.FetchUntreated()

	g.Response(http.StatusOK, e.Success, count)
}
