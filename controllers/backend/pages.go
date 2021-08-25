package backend

import (
	"eCommerce/pkg/e"
	"net/http"

	models "eCommerce/models/backend"

	"github.com/gin-gonic/gin"
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
