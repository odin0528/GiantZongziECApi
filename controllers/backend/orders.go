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

	if req.WithoutProducts != true {
		for index := range orders {
			orders[index].GetProducts()
		}
	}

	g.PaginationResponse(http.StatusOK, e.Success, orders, pagination)
}

func OrderStartPickup(c *gin.Context) {
	g := Gin{c}
	var query models.BatchOrderQuery
	err := c.BindJSON(&query)

	PlatformID, _ := c.Get("platform_id")
	AdminID, _ := c.Get("admin_id")

	err = DB.Model(models.Orders{}).Where(`
		id in ? AND 
		status = 21 AND 
		platform_id = ?
	`, query.IDs, PlatformID.(int)).Updates(map[string]interface{}{"status": 22, "picker_id": AdminID.(int)}).Error

	if err != nil {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}

	g.Response(http.StatusOK, e.Success, nil)
}

func OrderCancel(c *gin.Context) {
	g := Gin{c}
	var query models.OrderQuery
	err := c.BindJSON(&query)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	PlatformID, _ := c.Get("platform_id")

	affected := DB.Model(models.Orders{}).Where(`
		id = ? AND 
		(status = 21 AND payment = 2 OR status = 11 ) AND 
		platform_id = ?
	`, query.ID, PlatformID.(int)).Update("status", 99).RowsAffected

	if affected == 0 {
		g.Response(http.StatusOK, e.StatusNotFound, err)
		return
	}

	g.Response(http.StatusOK, e.Success, nil)
}

func OrderUntreated(c *gin.Context) {
	g := Gin{c}
	var query models.OrderQuery
	PlatformID, _ := c.Get("platform_id")
	query.PlatformID = PlatformID.(int)
	count := query.FetchUntreated()

	g.Response(http.StatusOK, e.Success, count)
}
