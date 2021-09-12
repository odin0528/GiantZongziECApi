package backend

import (
	"eCommerce/pkg/e"
	"net/http"

	. "eCommerce/internal/database"
	models "eCommerce/models/backend"

	"github.com/gin-gonic/gin"
)

func PromotionList(c *gin.Context) {
	g := Gin{c}
	var req models.PromotionListReq
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}
	PlatformID, _ := c.Get("platform_id")
	req.PlatformID = PlatformID.(int)
	promotions, pagination := req.FetchAll()

	g.PaginationResponse(http.StatusOK, e.Success, promotions, pagination)
}

func PromotionModify(c *gin.Context) {
	g := Gin{c}
	var req *models.Promotions
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	platformID, _ := c.Get("platform_id")
	req.PlatformID = platformID.(int)
	/* query := &models.CategoryQuery{
		ID: req.ID,
	}

	category := query.Fetch()
	if !category.Validate(platformID.(int)) {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	} */

	if req.StartTimestamp >= req.EndTimestamp {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	switch req.Type {
	case "sitewide-discount":
		if req.Mode == "total_qty" && req.Qty <= 0 {
			g.Response(http.StatusBadRequest, e.InvalidParams, err)
			return
		}

		if req.Mode == "total_price" && req.Money <= 0 {
			g.Response(http.StatusBadRequest, e.InvalidParams, err)
			return
		}

		if req.Method == "percent" && (req.Percent <= 0 || req.Percent >= 100) {
			g.Response(http.StatusBadRequest, e.InvalidParams, err)
			return
		}

		if req.Method == "discount" && req.Discount <= 0 {
			g.Response(http.StatusBadRequest, e.InvalidParams, err)
			return
		}
	}

	if req.ID == 0 {
		err = DB.Create(&req).Error
		if err != nil {
			g.Response(http.StatusBadRequest, e.InvalidParams, err)
			return
		}
	}

	g.Response(http.StatusOK, e.Success, nil)
}
