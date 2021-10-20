package backend

import (
	"eCommerce/pkg/e"
	"net/http"

	. "eCommerce/internal/database"
	models "eCommerce/models/backend"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func GetPagesList(c *gin.Context) {
	g := Gin{c}
	PlatformID, _ := c.Get("platform_id")
	req := &models.PageReq{
		PlatformID: PlatformID.(int),
	}

	pages, _ := req.GetPageList()

	g.Response(http.StatusOK, e.Success, pages)
}

func GetPageComponent(c *gin.Context) {
	g := Gin{c}
	var req models.PageReq
	err := c.ShouldBindUri(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	platformID, _ := c.Get("platform_id")
	req.PlatformID = platformID.(int)
	pages, err := req.Fetch()

	if pages.ID == 0 || err != nil {
		g.Response(http.StatusOK, e.StatusNotFound, nil)
		return
	}

	componentReq := models.PageComponentDraftQuery{
		PageID:     pages.ID,
		PlatformID: platformID.(int),
	}

	components := componentReq.FetchByPageID()
	for index, component := range components {
		componentQuery := models.PageComponentDataDraftQuery{
			ComID: component.ID,
		}
		componentData := componentQuery.FetchByComID()
		components[index].Data = componentData
	}

	g.Response(http.StatusOK, e.Success, components)

}

func PageRelease(c *gin.Context) {
	g := Gin{c}
	var req *models.PageReq
	c.BindJSON(&req)

	PlatformID, _ := c.Get("platform_id")
	req.PlatformID = PlatformID.(int)

	page, err := req.Fetch()

	if page.ID == 0 || err != nil {
		g.Response(http.StatusOK, e.StatusNotFound, nil)
		return
	}

	req.Clear()
	err = req.DeepDuplicate()

	if err != nil {
		g.Response(http.StatusInternalServerError, e.StatusInternalServerError, err)
		return
	}

	g.Response(http.StatusOK, e.Success, nil)
}

func PageModify(c *gin.Context) {
	g := Gin{c}
	var page models.Pages
	err := c.BindJSON(&page)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	validate := validator.New()
	err = validate.Struct(page)
	if err != nil {
		g.Response(http.StatusOK, e.InvalidParams, err.(validator.ValidationErrors)[0].Field())
		return
	}

	PlatformID, _ := c.Get("platform_id")
	page.PlatformID = PlatformID.(int)

	if page.ID == 0 {
		err = DB.Create(&page).Error
		if err != nil {
			g.Response(http.StatusBadRequest, e.StatusInternalServerError, err)
			return
		}

		g.Response(http.StatusOK, e.Success, nil)
		return
	} else {
		err = DB.Select("url", "title", "is_menu", "is_enabled").Updates(&page).Error
		if err != nil {
			g.Response(http.StatusBadRequest, e.StatusInternalServerError, err)
			return
		}

		g.Response(http.StatusOK, e.Success, nil)
		return
	}
}

/* func PageSort(c *gin.Context) {
	g := Gin{c}
	var pages []models.Pages
	err := c.BindJSON(&pages)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	for index, page := range pages {
		page.Sort = index
		err = DB.Select("sort").Updates(&page).Error
		if err != nil {
			g.Response(http.StatusBadRequest, e.StatusInternalServerError, err)
			return
		}
	}

	g.Response(http.StatusOK, e.Success, nil)
	return
} */
