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

func GetPage(c *gin.Context) {
	g := Gin{c}
	var req models.PageComponentDraft
	err := c.ShouldBindUri(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	components := req.FetchByPageID()
	g.Response(http.StatusOK, e.Success, components)

}
