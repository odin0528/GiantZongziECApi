package frontend

import (
	"crypto/sha256"
	models "eCommerce/models/frontend"
	"eCommerce/pkg/e"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"eCommerce/internal/ads"
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

func GuestCartsFetch(c *gin.Context) {
	g := Gin{c}
	query := &models.GuestCartsQuery{}
	err := c.BindJSON(&query)
	if err != nil {
		println(err.Error())
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	PlatformID, _ := c.Get("platform_id")
	query.PlatformID = PlatformID.(int)
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
	MemberEmail, _ := c.Get("member_email")
	p, _ := c.Get("platform")
	platform := p.(models.Platform)

	carts.PlatformID = PlatformID.(int)
	carts.MemberID = MemberID.(int)

	err = DB.Create(&carts).Error

	if platform.FBPixel != "" && platform.FBPixelToken != "" {
		em := sha256.Sum256([]byte(strings.ToLower(MemberEmail.(string))))
		params := ads.FbConversionParams{
			EventName:      "AddToCart",
			EventTime:      time.Now().Unix(),
			EventSourceUrl: fmt.Sprintf("https://%s/%s", platform.Hostname, fmt.Sprintf("products/%d", carts.ProductID)),
			ActionSource:   "website",
			UserData: ads.FbConversionUserData{
				EM:         hex.EncodeToString(em[:]),
				ExternalID: fmt.Sprintf("%d", carts.StyleID),
			},
		}

		ads.Send(params, platform.FBPixel, platform.FBPixelToken)
	}

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
		err = DB.Model(&models.Carts{}).Where("platform_id = ? AND member_id = ? AND product_id = ? AND style_id = ?", PlatformID.(int), MemberID.(int), product.ProductID, product.StyleID).
			Updates(map[string]interface{}{
				"buy_count": product.BuyCount,
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
