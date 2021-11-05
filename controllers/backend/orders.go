package backend

import (
	"eCommerce/pkg/e"
	"net/http"

	. "eCommerce/internal/database"
	"eCommerce/internal/money"
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
	PlatformID, _ := c.Get("platform_id")
	query.PlatformID = PlatformID.(int)
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
	PlatformID, _ := c.Get("platform_id")
	req.PlatformID = PlatformID.(int)
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
	PlatformID, _ := c.Get("platform_id")
	query.PlatformID = PlatformID.(int)
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

	err = DB.Select("status").Updates(&order).Error

	if err != nil {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}

	g.Response(http.StatusOK, e.Success, nil)
}

func OrderMakeShipmentNo(c *gin.Context) {
	g := Gin{c}
	var query models.OrderQuery
	err := c.BindJSON(&query)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}
	PlatformID, _ := c.Get("platform_id")
	query.PlatformID = PlatformID.(int)
	query.Status = 21
	order, err := query.Fetch()

	if err != nil {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}

	order.GetProducts()

	response, err := money.CreateLogisticsOrder(order)

	if err != nil {
		g.Response(http.StatusOK, e.StatusInternalServerError, err.Error())
		return
	}

	order.Status = 22
	order.LogisticsID = response.Get("AllPayLogisticsID")
	order.LogisticsStatus = 1
	order.LogisticsMsg = "託運單號建立完成"
	if order.Method == 1 {
		order.ShipmentNo = response.Get("BookingNote")
	} else {
		order.ShipmentNo = response.Get("ShipmentNo")
	}

	err = DB.Select("status", "logistics_id", "shipment_no", "logistics_status", "logistics_msg").Updates(&order).Error

	if err != nil {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}

	g.Response(http.StatusOK, e.Success, nil)
}

func OrderShipmentPrint(c *gin.Context) {
	g := Gin{c}
	var query models.BatchOrderQuery
	err := c.BindJSON(&query)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}
	PlatformID, _ := c.Get("platform_id")
	query.PlatformID = PlatformID.(int)
	query.Status = 22
	orders, err := query.FetchAll()

	if err != nil {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}

	ids := []string{}

	for _, order := range orders {
		ids = append(ids, order.LogisticsID)
	}

	g.Response(http.StatusOK, e.Success, money.GeneratePrintShipmentCheckMac(ids))
}

func OrderUntreated(c *gin.Context) {
	g := Gin{c}
	var query models.OrderQuery
	PlatformID, _ := c.Get("platform_id")
	query.PlatformID = PlatformID.(int)
	count := query.FetchUntreated()

	g.Response(http.StatusOK, e.Success, count)
}
