package backend

import (
	. "eCommerce/internal/database"
	"eCommerce/pkg/e"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PageReq struct {
	PageID     int `json:"page_id" uri:"page_id"`
	CustomerID int `json:"customer_id"`
}

type Pages struct {
	PageID     int    `json:"page_id"`
	CustomerID int    `json:"-"`
	Name       string `json:"name"`
	ReleasedAt int    `json:"released_at"`
	TimeDefault
}

func (req *PageReq) GetPageList() (pagesRowset []Pages, err error) {
	err = DB.Table("rel_customer_pages as rel").Select("rel.*, pages.name").Joins("inner join pages on rel.page_id = pages.id").
		Where("rel.customer_id = ?", req.CustomerID).
		Scan(&pagesRowset).Error
	return
}

func (req *PageReq) Fetch() (pages Pages) {
	DB.Table("rel_customer_pages").Where("page_id = ? and customer_id = ?", req.PageID, req.CustomerID).Scan(&pages)
	return
}

func (pages *Pages) Validate(customerID int, ctx gin.Context) {
	// data is not exist
	if pages.PageID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"http_status": http.StatusBadRequest,
			"code":        e.DataNotExist,
			"msg":         e.GetMsg(e.DataNotExist),
			"data":        nil,
		})
		ctx.Abort()
	}

	// The owner of the data is not the operator
	if pages.CustomerID != customerID {
		ctx.JSON(http.StatusForbidden, gin.H{
			"http_status": http.StatusForbidden,
			"code":        e.Forbidden,
			"msg":         e.GetMsg(e.Forbidden),
			"data":        nil,
		})
		ctx.Abort()
	}
}
