package backend

import (
	"net/http"

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

	component1 := models.PageComponentDraft{
		PageID: req.PageID,
		Sort:   req.Position1,
	}

	component2 := models.PageComponentDraft{
		PageID: req.PageID,
		Sort:   req.Position2,
	}

	component1.FetchBySort()
	component2.FetchBySort()
	component1.Sort = req.Position2
	component2.Sort = req.Position1
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

	component.Save()

	for _, data := range component.Data {
		data.PageID = component.PageID
		data.ComID = component.ID
		data.Save()
	}

	g.Response(http.StatusOK, e.Success, nil)
}
