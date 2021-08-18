package backend

import (
	"eCommerce/pkg/e"
	"net/http"
	"time"

	models "eCommerce/models/backend"

	"github.com/gin-gonic/gin"
)

func CategoryList(c *gin.Context) {
	g := Gin{c}
	var query models.CategoryQuery
	err := c.ShouldBindUri(&query)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}
	CustomerID, _ := c.Get("customer_id")
	query.CustomerID = CustomerID.(int)
	categories := query.FetchAll()

	var breadcrumbs []models.Category
	query.GetBreadcrumbs(&breadcrumbs)

	g.Response(http.StatusOK, e.Success, map[string]interface{}{"categories": categories, "breadcrumbs": breadcrumbs})
}

func CategoryChildList(c *gin.Context) {
	g := Gin{c}
	var query models.CategoryQuery
	err := c.ShouldBindUri(&query)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}
	CustomerID, _ := c.Get("customer_id")
	query.CustomerID = CustomerID.(int)
	categories := query.FetchAll()

	g.Response(http.StatusOK, e.Success, categories)
}

func CategoryCreate(c *gin.Context) {
	g := Gin{c}
	var category *models.Category
	err := c.BindJSON(&category)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	CustomerID, _ := c.Get("customer_id")

	query := &models.CategoryQuery{
		CustomerID: CustomerID.(int),
		ParentID:   category.ParentID,
	}

	category.CustomerID = CustomerID.(int)
	category.Sort = int(query.Count()) + 1
	category.CreatedAt = int(time.Now().Unix())
	category.UpdatedAt = int(time.Now().Unix())
	category.Create()

	g.Response(http.StatusOK, e.Success, category)
}

func CategoryModify(c *gin.Context) {
	g := Gin{c}
	var req *models.CategoryModifyReq
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	customerID, _ := c.Get("customer_id")
	query := &models.CategoryQuery{
		ID: req.ID,
	}

	category := query.Fetch()
	if !category.Validate(customerID.(int)) {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}

	// update data
	category.Title = req.Title
	category.UpdatedAt = int(time.Now().Unix())
	category.Update()

	g.Response(http.StatusOK, e.Success, nil)
}

func CategoryMove(c *gin.Context) {
	g := Gin{c}
	var req *models.CategoryMoveReq
	customerID, _ := c.Get("customer_id")
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	categoryQuery1 := models.CategoryQuery{
		ParentID: req.ParentID,
		Sort:     req.Sort,
	}

	categoryQuery2 := models.CategoryQuery{
		ParentID: req.ParentID,
		Sort:     req.Sort + req.Direction,
	}

	category1 := categoryQuery1.Fetch()
	category1.Validate(customerID.(int))
	if !category1.Validate(customerID.(int)) {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}
	category2 := categoryQuery2.Fetch()
	category2.Validate(customerID.(int))
	if !category2.Validate(customerID.(int)) {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}
	category2.Sort = req.Sort
	category1.Sort = req.Sort + req.Direction
	category1.Update()
	category2.Update()

	g.Response(http.StatusOK, e.Success, nil)
}

func CategoryDelete(c *gin.Context) {
	g := Gin{c}
	var req *models.Category
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	customerID, _ := c.Get("customer_id")
	query := &models.CategoryQuery{
		ID: req.ID,
	}

	category := query.Fetch()
	if !category.Validate(customerID.(int)) {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}

	category.Delete()

	g.Response(http.StatusOK, e.Success, nil)
}
