package frontend

import (
	"eCommerce/pkg/e"
	"net/http"

	models "eCommerce/models/frontend"

	"github.com/gin-gonic/gin"
)

func GetUserOrders(c *gin.Context) {
	g := Gin{c}
	var req *models.OrderQuery
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}
	MemeberID, _ := c.Get("member_id")
	req.MemberID = MemeberID.(int)

	if req.MemberID == 0 {
		g.Response(http.StatusUnauthorized, e.Unauthorized, err)
		return
	}

	orders, pagination := req.FetchAll()

	g.PaginationResponse(http.StatusOK, e.Success, orders, pagination)
}
