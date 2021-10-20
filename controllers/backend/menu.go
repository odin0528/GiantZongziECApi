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
	var query models.CategoryQuery
	err := c.ShouldBindUri(&query)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}
	PlatformID, _ := c.Get("platform_id")
	query.PlatformID = PlatformID.(int)
	categories := query.FetchAll()

	var breadcrumbs []models.Category
	query.GetBreadcrumbs(&breadcrumbs)

	g.Response(http.StatusOK, e.Success, map[string]interface{}{"categories": categories, "breadcrumbs": breadcrumbs})
}

func MenuModify(c *gin.Context) {
	g := Gin{c}
	var req *models.Menu
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	platformID, _ := c.Get("platform_id")

	if req.ID == 0 {
		req.PlatformID = platformID.(int)
		DB.Debug().Create(&req)
	} else {
		err = DB.Debug().Select("title", "link", "link_type", "is_enabled").Where("id = ? and platform_id = ?", req.ID, platformID.(int)).Updates(&req).Error
	}

	g.Response(http.StatusOK, e.Success, nil)
}

func MenuMove(c *gin.Context) {
	g := Gin{c}
	var req *models.CategoryMoveReq
	platformID, _ := c.Get("platform_id")
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	categoryQuery1 := models.CategoryQuery{
		ParentID:   req.ParentID,
		PlatformID: platformID.(int),
		Sort:       req.Sort,
	}

	categoryQuery2 := models.CategoryQuery{
		ParentID:   req.ParentID,
		PlatformID: platformID.(int),
		Sort:       req.Sort + req.Direction,
	}

	category1 := categoryQuery1.Fetch()
	if category1.ID == 0 {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}
	category2 := categoryQuery2.Fetch()
	if category2.ID == 0 {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}
	category2.Sort = req.Sort
	category1.Sort = req.Sort + req.Direction
	category1.Update()
	category2.Update()

	g.Response(http.StatusOK, e.Success, nil)
}

func MenuDelete(c *gin.Context) {
	g := Gin{c}
	var req *models.Category
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	platformID, _ := c.Get("platform_id")
	query := &models.CategoryQuery{
		ID: req.ID,
	}

	category := query.Fetch()
	if !category.Validate(platformID.(int)) {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}

	category.Delete()

	g.Response(http.StatusOK, e.Success, nil)
}
