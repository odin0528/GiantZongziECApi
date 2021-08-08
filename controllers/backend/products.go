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

func ProductFetch(c *gin.Context) {
	g := Gin{c}
	var query models.ProductQuery
	err := c.ShouldBindUri(&query)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}
	CustomerID, _ := c.Get("customer_id")
	query.CustomerID = CustomerID.(int)
	product := query.Fetch()
	product.GetPhotos()
	product.GetStyle()
	product.GetSubStyle()
	product.GetStyleTable()

	g.Response(http.StatusOK, e.Success, product)
}

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

		sort := 1
		for _, photo := range req.Photos {
			//有找到base64的編碼關鍵字
			if strings.Index(photo.Img, ",") > 0 {
				filename := fmt.Sprintf("/upload/%08d/products/%08d/%d", customerID.(int), req.ID, time.Now().UnixNano())
				productPhotos := models.ProductPhotos{
					ProductID:  req.ID,
					CustomerID: customerID.(int),
					Img:        uploader.Thumbnail(filename, photo.Img, 720),
					Sort:       sort,
				}
				productPhotos.Create()
				sort++
			}
		}

		for index, style := range req.Style {
			style.ProductID = req.ID
			style.CustomerID = customerID.(int)
			style.Sort = index
			if strings.Index(style.Img, ",") > 0 {
				filename := fmt.Sprintf("/upload/%08d/products/%08d/%d", customerID.(int), req.ID, time.Now().UnixNano())
				style.Img = uploader.Thumbnail(filename, style.Img, 720)
			}
			style.Create()
		}

		for index, style := range req.SubStyle {
			style.ProductID = req.ID
			style.CustomerID = customerID.(int)
			style.Sort = index
			style.Create()
		}

		for index, list := range req.StyleTable {
			for _, item := range list {
				item.ProductID = req.ID
				item.CustomerID = customerID.(int)
				item.Group = index
				item.Create()
			}
		}
	} else {
		// 編輯產品
		var query models.ProductQuery
		query.ID = req.ID
		product := query.Fetch()
		if !product.Validate(customerID.(int)) {
			g.Response(http.StatusBadRequest, e.DataNotExist, err)
			return
		}
		req.Update()

		sort := 1
		for _, photo := range req.Photos {
			// 如果沒有id 新增照片
			if photo.ID == 0 {
				if strings.Index(photo.Img, ",") > 0 {
					filename := fmt.Sprintf("/upload/%08d/products/%08d/%d", customerID.(int), req.ID, time.Now().UnixNano())
					productPhotos := models.ProductPhotos{
						ProductID:  req.ID,
						CustomerID: customerID.(int),
						Img:        uploader.Thumbnail(filename, photo.Img, 720),
						Sort:       sort,
					}
					productPhotos.Create()
					sort++
				}
			} else {
				// 如果是空值，就當它刪除了
				if photo.Img == "" {
					productPhotos := models.ProductPhotos{
						ID:         photo.ID,
						CustomerID: customerID.(int),
					}
					productPhotos.Fetch()
					uploader.DeletePhoto(productPhotos.Img)
					productPhotos.Delete()
				} else {
					productPhotos := models.ProductPhotos{
						ID:   photo.ID,
						Img:  photo.Img,
						Sort: sort,
					}
					//有找到base64的編碼關鍵字
					if strings.Index(photo.Img, ",") > 0 {
						filename := fmt.Sprintf("/upload/%08d/products/%08d/%d", customerID.(int), req.ID, time.Now().UnixNano())
						productPhotos.Img = uploader.Thumbnail(filename, photo.Img, 720)
					} else {
						productPhotos.Img = photo.Img
					}

					productPhotos.Update()
					sort++
				}
			}
		}

		deleteIds := []int{}
		for index, style := range req.Style {
			style.Sort = index
			if strings.Index(style.Img, ",") > 0 {
				filename := fmt.Sprintf("/upload/%08d/products/%08d/%d", customerID.(int), req.ID, time.Now().UnixNano())
				style.Img = uploader.Thumbnail(filename, style.Img, 720)
			}
			if style.ID == 0 {
				style.ProductID = req.ID
				style.CustomerID = customerID.(int)
				style.Create()
			} else {
				style.Update()
			}
			deleteIds = append(deleteIds, style.ID)
		}

		productStyle := &models.ProductStyle{
			ProductID:  req.ID,
			CustomerID: customerID.(int),
		}
		productStyle.DeleteNotExistStyle(deleteIds)

		deleteIds = []int{}
		for index, style := range req.SubStyle {
			style.Sort = index

			if style.ID == 0 {
				style.ProductID = req.ID
				style.CustomerID = customerID.(int)
				style.Create()
			} else {
				style.Update()
			}
			deleteIds = append(deleteIds, style.ID)
		}

		productSubStyle := &models.ProductSubStyle{
			ProductID:  req.ID,
			CustomerID: customerID.(int),
		}
		productSubStyle.DeleteNotExistStyle(deleteIds)

		deleteIds = []int{}
		for index, list := range req.StyleTable {
			for _, item := range list {
				item.ProductID = req.ID
				item.CustomerID = customerID.(int)
				item.Group = index
				if item.ID == 0 {
					item.Create()
				} else {
					item.Update()
				}
				deleteIds = append(deleteIds, item.ID)
			}
		}

		if len(deleteIds) > 0 {
			styleTable := &models.ProductStyleTable{
				ProductID:  req.ID,
				CustomerID: customerID.(int),
			}
			styleTable.DeleteNotExistStyle(deleteIds)
		}
	}

	g.Response(http.StatusOK, e.Success, nil)
}
