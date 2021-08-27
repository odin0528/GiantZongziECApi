package backend

import (
	models "eCommerce/models/backend"
	"eCommerce/pkg/e"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MemberList(c *gin.Context) {
	g := Gin{c}
	var query models.MemberQuery
	err := c.BindJSON(&query)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}
	PlatformID, _ := c.Get("platform_id")
	query.PlatformID = PlatformID.(int)
	members, pager := query.FetchAll()

	g.PaginationResponse(http.StatusOK, e.Success, members, pager)
}
