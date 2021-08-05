package backend

import (
	"eCommerce/pkg/e"
	"net/http"

	models "eCommerce/models/backend"

	"github.com/gin-gonic/gin"
	"github.com/liudng/godump"
)

func ProductModify(c *gin.Context) {
	g := Gin{c}
	var req *models.Products
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	godump.Dump(req)

	if req.ID == 0 {
		req.Create()

		for _, list := range req.StyleTable {
			for _, item := range list {
				item.ProductID = req.ID
				item.Create()
			}
		}
	}

	g.Response(http.StatusOK, e.Success, nil)
}
