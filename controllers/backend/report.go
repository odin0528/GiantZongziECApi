package backend

import (
	"eCommerce/pkg/e"
	"net/http"

	models "eCommerce/models/backend"

	"github.com/gin-gonic/gin"
)

func StockReport(c *gin.Context) {
	g := Gin{c}
	var req models.ReportType
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	var pagination models.Pagination
	var products []models.Products

	PlatformID, _ := c.Get("platform_id")
	req.PlatformID = PlatformID.(int)

	switch req.Type {
	case "low_stock":
		products, pagination = req.LowStock()
		for index := range products {
			products[index].GetLowStockStyleTable()
		}
	case "over_sale":
		products, pagination = req.OverSale()
		for index := range products {
			products[index].GetOverSaleStyleTable()
		}
	case "wait_for_delivery":
		products, pagination = req.WaitForDelivery()
		for index := range products {
			products[index].GetWaitDeliveryStyleTable()
		}
	}

	g.PaginationResponse(http.StatusOK, e.Success, products, pagination)
}