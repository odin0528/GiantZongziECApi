package backend

import (
	"net/http"
	"time"

	"eCommerce/pkg/e"

	models "eCommerce/models/backend"

	"github.com/gin-gonic/gin"
)

func DraftComponentDelete(c *gin.Context) {
	g := Gin{c}
	var component *models.PageComponentDraft
	err := c.BindJSON(&component)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	component.FetchBySort()
	component.Delete()
	component.DeleteChildren()

	g.Response(http.StatusOK, e.Success, nil)
}

func DraftComponentEdit(c *gin.Context) {
	g := Gin{c}
	var component *models.PageComponentDraft
	err := c.BindJSON(&component)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	component.FetchBySort()
	component.DeleteChildren()

	for _, data := range component.Data {
		data.PageID = component.PageID
		data.ComID = component.ID
		/*i := strings.Index(data.Img, ",") //有找到base64的編碼關鍵字
		if i > 0 {
			filename := fmt.Sprintf("updata/dmo/cms/%d_%d_%d.jpg", data.PageID, data.ComID, index)
			blob, _ := base64.StdEncoding.DecodeString(data.Img[i+1:])
			data.Img = common.Storage(filename, blob)
		}*/

		data.Save()
	}

	g.Response(http.StatusOK, e.Success, nil)
}

func DraftComponentChange(c *gin.Context) {
	g := Gin{c}
	var req *models.PageComponentDraftChangeReq
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	component := models.PageComponentDraft{
		PageID: req.PageID,
		Sort:   req.Position1,
	}

	component.FetchBySort()
	component.Sort = req.Position2
	component.Update()

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
	CustomerID, _ := c.Get("customer_id")
	component.CustomerID = CustomerID.(int)
	component.CreatedAt = int(time.Now().Unix())
	component.UpdatedAt = int(time.Now().Unix())
	component.Save()

	for _, data := range component.Data {
		data.PageID = component.PageID
		data.ComID = component.ID
		data.Save()
	}

	g.Response(http.StatusOK, e.Success, nil)
}
