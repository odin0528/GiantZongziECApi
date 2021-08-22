package backend

import (
	. "eCommerce/internal/database"
	"eCommerce/pkg/e"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PageReq struct {
	PageID     int `json:"page_id" uri:"page_id"`
	PlatformID int `json:"platform_id"`
}

type Pages struct {
	PageID     int    `json:"page_id"`
	PlatformID int    `json:"-"`
	Name       string `json:"name"`
	ReleasedAt int    `json:"released_at"`
	TimeDefault
}

func (req *PageReq) GetPageList() (pagesRowset []Pages, err error) {
	err = DB.Table("rel_platform_pages as rel").Select("rel.*, pages.name").Joins("inner join pages on rel.page_id = pages.id").
		Where("rel.platform_id = ?", req.PlatformID).
		Scan(&pagesRowset).Error
	return
}

func (req *PageReq) Fetch() (pages Pages) {
	DB.Table("rel_platform_pages").Where("page_id = ? and platform_id = ?", req.PageID, req.PlatformID).Scan(&pages)
	return
}

func (pages *Pages) Validate(platformID int, ctx gin.Context) {
	// data is not exist
	if pages.PageID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"http_status": http.StatusBadRequest,
			"code":        e.StatusNotFound,
			"msg":         e.GetMsg(e.StatusNotFound),
			"data":        nil,
		})
		ctx.Abort()
	}

	// The owner of the data is not the operator
	if pages.PlatformID != platformID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"http_status": http.StatusForbidden,
			"code":        e.Forbidden,
			"msg":         e.GetMsg(e.Forbidden),
			"data":        nil,
		})
		ctx.Abort()
	}
}
