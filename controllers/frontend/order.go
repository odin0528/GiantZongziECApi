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
	switch order.Payment {
	case 1:
		order.Status = 11
	case 2:
		order.Status = 21
	}
	order.Create()

	for _, product := range order.Products {
		for _, style := range product.Styles {
			orderProduct := models.OrderProducts{
				OrderID:    order.ID,
				ProductID:  product.ProductID,
				StyleID:    style.StyleID,
				Qty:        style.Qty,
				Price:      style.Price,
				Total:      float32(style.Qty) * style.Price,
				Title:      product.Title,
				StyleTitle: style.Title + style.SubTitle,
				Photo:      style.Photo,
				Sku:        style.Sku,
			}

			orderProduct.Create()
		}
	}

	g.Response(http.StatusOK, e.Success, nil)
}
