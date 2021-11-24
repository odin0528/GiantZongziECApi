package frontend

import (
	models "eCommerce/models/frontend"
	"eCommerce/pkg/e"
	"net/http"

	"github.com/gin-gonic/gin"

	. "eCommerce/internal/database"
)

func CartsFetch(c *gin.Context) {
	g := Gin{c}
	query := &models.CartsQuery{}

	PlatformID, _ := c.Get("platform_id")
	MemberID, _ := c.Get("member_id")

	query.PlatformID = PlatformID.(int)
	query.MemberID = MemberID.(int)

	carts, err := query.FetchAll()

	if err != nil {
		g.Response(http.StatusBadRequest, e.StatusInternalServerError, err)
		return
	}

	g.Response(http.StatusOK, e.Success, carts)
}

func CartsAddProduct(c *gin.Context) {
	g := Gin{c}
	var carts *models.Carts
	err := c.BindJSON(&carts)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	PlatformID, _ := c.Get("platform_id")
	MemberID, _ := c.Get("member_id")

	carts.PlatformID = PlatformID.(int)
	carts.MemberID = MemberID.(int)

	err = DB.Create(&carts).Error

	if err != nil {
		g.Response(http.StatusBadRequest, e.StatusInternalServerError, err)
		return
	}

	g.Response(http.StatusOK, e.Success, nil)
}

func CartsResetProduct(c *gin.Context) {
	g := Gin{c}
	var carts []models.Carts
	err := c.BindJSON(&carts)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	PlatformID, _ := c.Get("platform_id")
	MemberID, _ := c.Get("member_id")

	/* memberCarts := models.Carts{
		MemberID:   MemberID.(int),
		PlatformID: PlatformID.(int),
	}
	memberCarts.Clean() */

	for _, product := range carts {
		err = DB.Debug().Model(&models.Carts{}).Where("platform_id = ? AND member_id = ? AND product_id = ? AND style_id = ?", PlatformID.(int), MemberID.(int), product.ProductID, product.StyleID).
			Updates(map[string]interface{}{
				"qty": product.Qty,
			}).Error

		if err != nil {
			g.Response(http.StatusBadRequest, e.StatusInternalServerError, err)
			return
		}
	}

	g.Response(http.StatusOK, e.Success, nil)
}

func CartsUpdate(c *gin.Context) {
	g := Gin{c}
	var carts *models.Carts
	err := c.BindJSON(&carts)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	PlatformID, _ := c.Get("platform_id")
	MemberID, _ := c.Get("member_id")

	carts.PlatformID = PlatformID.(int)
	carts.MemberID = MemberID.(int)

	err = carts.Update()

	if err != nil {
		g.Response(http.StatusBadRequest, e.StatusInternalServerError, err)
		return
	}

	g.Response(http.StatusOK, e.Success, nil)
}

func CartsRemoveProduct(c *gin.Context) {
	g := Gin{c}
	var carts *models.Carts
	err := c.BindJSON(&carts)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	PlatformID, _ := c.Get("platform_id")
	MemberID, _ := c.Get("member_id")

	carts.PlatformID = PlatformID.(int)
	carts.MemberID = MemberID.(int)

	err = carts.Delete()

	if err != nil {
		g.Response(http.StatusBadRequest, e.StatusInternalServerError, err)
		return
	}

	g.Response(http.StatusOK, e.Success, nil)
}
