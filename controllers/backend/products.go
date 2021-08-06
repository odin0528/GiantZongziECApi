package backend

import (
	"eCommerce/pkg/e"
	"fmt"
	"net/http"
	"strings"
	"time"

	"eCommerce/internal/uploader"
	models "eCommerce/models/backend"

	"github.com/gin-gonic/gin"
)

func ProductModify(c *gin.Context) {
	g := Gin{c}
	var req *models.Products
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}
	customerID, _ := c.Get("customer_id")

	if req.ID == 0 {
		req.CustomerID = customerID.(int)
		req.Create()

		for _, photo := range req.Photos {
			productPhotos := models.ProductPhotos{
				ProductID: req.ID,
			}

			i := strings.Index(photo.Img, ",") //有找到base64的編碼關鍵字
			if i > 0 {
				filename := fmt.Sprintf("/upload/%08d/products/%08d/%d", customerID.(int), req.ID, time.Now().UnixNano())
				productPhotos.Img = uploader.Thumbnail(filename, photo.Img, 720)
				productPhotos.Create()
			}
		}

		for index, list := range req.StyleTable {
			for key, item := range list {
				item.ProductID = req.ID
				item.Create()

				if key == 0 {
					i := strings.Index(req.StylePhotos[index].Img, ",") //有找到base64的編碼關鍵字
					if i > 0 {
						productStylePhotos := models.ProductStylePhotos{
							ProductID:      req.ID,
							ProductStyleID: item.ID,
						}
						filename := fmt.Sprintf("/upload/%08d/products/%08d/%d", customerID.(int), req.ID, time.Now().UnixNano())
						productStylePhotos.Img = uploader.Thumbnail(filename, req.StylePhotos[index].Img, 720)
						productStylePhotos.Create()
					}
				}
			}
		}
	}

	g.Response(http.StatusOK, e.Success, nil)
}
