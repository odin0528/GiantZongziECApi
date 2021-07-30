package backend

import (
	"eCommerce/pkg/e"
	"net/http"
	"time"

	models "eCommerce/models/backend"

	"github.com/gin-gonic/gin"
)

func CategoryCreate(c *gin.Context) {
	g := Gin{c}
	var category *models.Category
	err := c.BindJSON(&category)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}
	CustomerID, _ := c.Get("customer_id")
	category.CustomerID = CustomerID.(int)
	category.CreatedAt = int(time.Now().Unix())
	category.UpdatedAt = int(time.Now().Unix())
	category.Create()

	g.Response(http.StatusOK, e.Success, nil)
}
