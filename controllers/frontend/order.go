package frontend

import (
	models "eCommerce/models/frontend"
	"eCommerce/pkg/e"
	"net/http"

	"github.com/gin-gonic/gin"
)

func OrderCreate(c *gin.Context) {
	g := Gin{c}
	var order *models.OrderCreateRequest
	err := c.BindJSON(&order)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	CustomerID, _ := c.Get("customer_id")

	order.CustomerID = CustomerID.(int)
	order.Create()

	g.Response(http.StatusOK, e.Success, order)
}
