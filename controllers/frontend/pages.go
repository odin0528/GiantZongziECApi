package frontend

import (
	"eCommerce/pkg/e"
	"net/http"

	models "eCommerce/models/frontend"

	"github.com/gin-gonic/gin"
)

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
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}

	componentReq := models.PageComponentQuery{
		PageID:     pages.ID,
		PlatformID: platformID.(int),
	}

	components := componentReq.FetchByPageID()
	for index, component := range components {
		componentQuery := models.PageComponentDataQuery{
			ComID: component.ID,
		}
		componentData := componentQuery.FetchByComID()
		components[index].Data = componentData
	}

	g.Response(http.StatusOK, e.Success, components)

}
