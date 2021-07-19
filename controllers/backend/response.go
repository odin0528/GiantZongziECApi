package backend

import (
	"eCommerce/pkg/e"

	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.JSON(httpCode, gin.H{
		"http_status": httpCode,
		"code":        errCode,
		"msg":         e.GetMsg(errCode),
		"data":        data,
	})
}

/*func (g *Gin) PaginationResponse(httpCode, errCode int, data interface{}, pager interface{}) {
	g.C.JSON(httpCode, gin.H{
		"http_status": httpCode,
		"code":        errCode,
		"msg":         e.GetMsg(errCode),
		"data":        data,
		"pager":       pager,
	})
}*/
