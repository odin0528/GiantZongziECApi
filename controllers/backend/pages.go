package backend

import (
	"eCommerce/pkg/e"
	"net/http"

	models "eCommerce/models/backend"

	"github.com/gin-gonic/gin"
)

func GetPagesList(c *gin.Context) {
	g := Gin{c}
	req := &models.PageReq{
		CustomerID: 1,
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

	customerID, _ := c.Get("customer_id")
	req.CustomerID = customerID.(int)
	pages := req.Fetch()
	pages.Validate(customerID.(int), *c)

	componentReq := models.PageComponentDraftQuery{
		PageID:     pages.PageID,
		CustomerID: customerID.(int),
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
