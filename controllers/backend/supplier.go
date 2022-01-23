package backend

import (
	"eCommerce/pkg/e"
	"fmt"
	"net/http"

	. "eCommerce/internal/database"
	"eCommerce/internal/rdb"
	models "eCommerce/models/backend"

	"github.com/gin-gonic/gin"
)

func SupplierList(c *gin.Context) {
	g := Gin{c}
	var query models.SupplierQuery
	PlatformID, _ := c.Get("platform_id")
	query.PlatformID = PlatformID.(int)
	menus, err := query.FetchAll()

	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	g.Response(http.StatusOK, e.Success, menus)
}

func SupplierModify(c *gin.Context) {
	g := Gin{c}
	var req *models.Supplier
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	platformID, _ := c.Get("platform_id")
	req.PlatformID = platformID.(int)

	if req.ID == 0 {
		err = DB.Create(&req).Error
		if err != nil {
			g.Response(http.StatusBadRequest, e.InvalidParams, err)
			return
		}
	} else {
		err = DB.Updates(&req).Error
		if err != nil {
			g.Response(http.StatusBadRequest, e.InvalidParams, err)
			return
		}
	}
	rdb.Del(fmt.Sprintf("platform_%d_supplier", req.PlatformID))

	g.Response(http.StatusOK, e.Success, nil)
}

func SupplierDelete(c *gin.Context) {
	g := Gin{c}
	var req *models.Supplier
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	platformID, _ := c.Get("platform_id")

	err = DB.Debug().Where("platform_id = ?", platformID.(int)).Delete(&req).Error
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	rdb.Del(fmt.Sprintf("platform_%d_supplier", platformID.(int)))

	g.Response(http.StatusOK, e.Success, nil)
}
