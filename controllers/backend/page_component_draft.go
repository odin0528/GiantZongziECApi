package backend

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"strings"
	"time"

	. "eCommerce/internal/uploader"
	"eCommerce/pkg/e"

	models "eCommerce/models/backend"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
)

func DraftComponentDelete(c *gin.Context) {
	g := Gin{c}
	var req *models.PageComponentDraftQuery
	platformID, _ := c.Get("platform_id")
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	req.PlatformID = platformID.(int)

	component := req.Fetch()
	component.Validate(platformID.(int))
	if !component.Validate(platformID.(int)) {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}

	component.Delete()
	component.DeleteChildren()

	g.Response(http.StatusOK, e.Success, nil)
}

func DraftComponentEdit(c *gin.Context) {
	g := Gin{c}
	var req *models.PageComponentDraftEditReq
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	platformID, _ := c.Get("platform_id")
	componentQurey := &models.PageComponentDraftQuery{
		PageID:     req.PageID,
		PlatformID: platformID.(int),
		Sort:       req.Sort,
	}

	component := componentQurey.Fetch()
	if !component.Validate(platformID.(int)) {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}

	// update data
	component.Title = req.Data.Title
	component.Text = req.Data.Text
	component.UpdatedAt = int(time.Now().Unix())
	component.Update()

	component.DeleteChildren()
	for _, data := range req.Data.Data {
		data.PageID = component.PageID
		data.ComID = component.ID
		i := strings.Index(data.Img, ",") //有找到base64的編碼關鍵字
		if i > 0 {
			filename := fmt.Sprintf("/upload/%08d/%08d/%08d", platformID.(int), data.PageID, data.ComID)
			blob, _ := base64.StdEncoding.DecodeString(data.Img[i+1:])
			var img image.Image
			var maxSize uint
			buff := new(bytes.Buffer)

			switch component.Type {
			case "image":
				maxSize = 720
			default:
				maxSize = 1100
			}

			switch strings.TrimSuffix(data.Img[5:i], ";base64") {
			case "image/png":
				img, _ = png.Decode(bytes.NewReader(blob))
				thumbnail := resize.Resize(maxSize, 0, img, resize.Lanczos3)
				png.Encode(buff, thumbnail)
			case "image/jpeg":
				img, _ = jpeg.Decode(bytes.NewReader(blob))
				thumbnail := resize.Resize(maxSize, 0, img, resize.Lanczos3)
				jpeg.Encode(buff, thumbnail, nil)
			}

			result, _ := Uploader.Upload(&s3manager.UploadInput{
				Bucket: aws.String(os.Getenv("AWS_BUCKET")),
				Key:    aws.String(filename),
				Body:   bytes.NewReader(buff.Bytes()),
			})

			data.Img = result.Location
		}

		data.Save()
	}

	g.Response(http.StatusOK, e.Success, nil)
}

func DraftComponentChange(c *gin.Context) {
	g := Gin{c}
	var req *models.PageComponentDraftChangeReq
	platformID, _ := c.Get("platform_id")
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	componentQuery1 := models.PageComponentDraftQuery{
		PageID:     req.PageID,
		Sort:       req.Sort,
		PlatformID: platformID.(int),
	}

	componentQuery2 := models.PageComponentDraftQuery{
		PageID:     req.PageID,
		Sort:       req.NewSort,
		PlatformID: platformID.(int),
	}

	component1 := componentQuery1.Fetch()
	if !component1.Validate(platformID.(int)) {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}
	component2 := componentQuery2.Fetch()
	if !component2.Validate(platformID.(int)) {
		g.Response(http.StatusBadRequest, e.StatusNotFound, err)
		return
	}
	component1.Sort = req.NewSort
	component2.Sort = req.Sort
	component1.Update()
	component2.Update()

	g.Response(http.StatusOK, e.Success, nil)
}

func DraftComponentCreate(c *gin.Context) {
	g := Gin{c}
	var component *models.PageComponentDraft
	err := c.BindJSON(&component)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}
	PlatformID, _ := c.Get("platform_id")
	component.PlatformID = PlatformID.(int)
	component.CreatedAt = int(time.Now().Unix())
	component.UpdatedAt = int(time.Now().Unix())
	component.Create()

	for _, data := range component.Data {
		data.PageID = component.PageID
		data.ComID = component.ID
		data.Save()
	}

	g.Response(http.StatusOK, e.Success, component)
}
