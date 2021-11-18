package backend

import (
	"eCommerce/pkg/e"
	"net/http"

	. "eCommerce/internal/database"
	models "eCommerce/models/backend"

	"github.com/gin-gonic/gin"
)

func MenuList(c *gin.Context) {
	g := Gin{c}
	var query models.MenuQuery
	PlatformID, _ := c.Get("platform_id")
	query.PlatformID = PlatformID.(int)
	menus, err := query.FetchAll()

	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	g.Response(http.StatusOK, e.Success, menus)
}

func MenuModify(c *gin.Context) {
	g := Gin{c}
	var req *models.Menus
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	platformID, _ := c.Get("platform_id")

	if req.ID == 0 {
		req.PlatformID = platformID.(int)
		DB.Create(&req)
	} else {
		err = DB.Select("title", "link", "link_type", "is_enabled").Where("id = ? and platform_id = ?", req.ID, platformID.(int)).Updates(&req).Error
	}

	g.Response(http.StatusOK, e.Success, req.ID)
}

func MenuMove(c *gin.Context) {
	g := Gin{c}
	var req *models.MenuMoveReq
	platformID, _ := c.Get("platform_id")
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	menuQuery1 := models.MenuQuery{
		PlatformID: platformID.(int),
		Sort:       req.Sort,
	}

	menuQuery2 := models.MenuQuery{
		PlatformID: platformID.(int),
		Sort:       req.Sort + req.Direction,
	}

	menu1 := menuQuery1.Fetch()
	if menu1.ID == 0 {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}
	menu2 := menuQuery2.Fetch()
	if menu2.ID == 0 {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}
	menu2.Sort = req.Sort
	menu1.Sort = req.Sort + req.Direction

	DB.Save(&menu1)
	DB.Save(&menu2)

	g.Response(http.StatusOK, e.Success, nil)
}

func MenuDelete(c *gin.Context) {
	g := Gin{c}
	var req *models.Menus
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	platformID, _ := c.Get("platform_id")
	query := &models.MenuQuery{
		ID:         req.ID,
		PlatformID: platformID.(int),
	}

	menu := query.Fetch()
	menu.Delete()

	g.Response(http.StatusOK, e.Success, nil)
}
