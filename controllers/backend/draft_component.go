package backend

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"eCommerce/pkg/e"

	models "eCommerce/models/backend"

	"github.com/gin-gonic/gin"
)

func DraftComponentDelete(c *gin.Context) {
	g := Gin{c}
	var component *models.DraftComponent
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
	var component *models.DraftComponent
	err := c.BindJSON(&component)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	component.FetchBySort()
	component.DeleteChildren()

	component.Data.DmoID = component.DmoID
	component.Data.ComID = component.ID
	component.Data.Save()

	if component.Type != "place" &&
		component.Type != "event" &&
		component.Type != "store" &&
		component.Type != "food" &&
		component.Type != "play" &&
		component.Type != "attraction" {

		for index, list := range component.Data.List {
			list.DmoID = component.DmoID
			list.ComID = component.ID
			i := strings.Index(list.Img, ",") //有找到base64的編碼關鍵字
			if i > 0 {
				filename := fmt.Sprintf("updata/dmo/cms/%d_%d_%d.jpg", list.DmoID, list.ComID, index)
				blob, _ := base64.StdEncoding.DecodeString(list.Img[i+1:])
				list.Img = common.Storage(filename, blob)
			}

			list.Save()
		}
	}

	g.Response(http.StatusOK, e.Success, nil)
}

func DraftComponentChange(c *gin.Context) {
	g := Gin{c}
	var req *models.DraftComponentChangeReq
	err := c.BindJSON(&req)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	component1 := models.DraftComponent{
		DmoID: req.DmoID,
		Sort:  req.Position1,
	}

	component2 := models.DraftComponent{
		DmoID: req.DmoID,
		Sort:  req.Position2,
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
	var component *models.DraftComponent
	err := c.BindJSON(&component)
	if err != nil {
		g.Response(http.StatusBadRequest, e.InvalidParams, err)
		return
	}

	component.Save()
	component.Data.DmoID = component.DmoID
	component.Data.ComID = component.ID
	component.Data.Save()

	if component.Type != "place" &&
		component.Type != "event" &&
		component.Type != "store" &&
		component.Type != "food" &&
		component.Type != "play" &&
		component.Type != "attraction" {

		for _, list := range component.Data.List {
			list.DmoID = component.DmoID
			list.ComID = component.ID
			list.Save()
		}
	}

	g.Response(http.StatusOK, e.Success, nil)
}
