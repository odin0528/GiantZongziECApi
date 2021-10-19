package backend

import (
	. "eCommerce/internal/database"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type PageReq struct {
	ID         int `json:"id" uri:"id"`
	PlatformID int `json:"platform_id"`
}

type Pages struct {
	ID         int    `json:"id"`
	PlatformID int    `json:"-"`
	Url        string `json:"url"`
	Title      string `json:"title"`
	IsMenu     bool   `json:"is_menu"`
	IsEnabled  bool   `json:"is_enabled"`
	Sort       int    `json:"sort"`
	ReleasedAt int    `json:"released_at"`
	DeletedAt  soft_delete.DeletedAt
	TimeDefault
}

func (req *PageReq) GetPageList() (pages []Pages, err error) {
	err = DB.Model(&Pages{}).
		Where("platform_id = ?", req.PlatformID).Order("sort ASC").
		Scan(&pages).Error
	return
}

func (req *PageReq) Fetch() (pages Pages, err error) {
	err = DB.Model(&Pages{}).Where("id = ? and platform_id = ?", req.ID, req.PlatformID).Scan(&pages).Error
	return
}

func (req *PageReq) Clear() {
	DB.Exec("DELETE FROM page_component WHERE page_id = ?", req.ID)
	DB.Exec("DELETE FROM page_component_data WHERE page_id = ?", req.ID)
}

func (req *PageReq) DeepDuplicate() error {
	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`
			INSERT INTO page_component 
			(id, platform_id, page_id, sort, component_name, type, title, text, created_at, updated_at) 
			(SELECT id, platform_id, page_id, sort, component_name, type, title, text, created_at, updated_at FROM page_component_draft WHERE page_id = ?)
		`, req.ID, req.PlatformID).Error; err != nil {
			return err
		}
		if err := tx.Exec(`
			INSERT INTO page_component_data (id, page_id, com_id, title, img, link_type, link, text, start_time, end_time) 
			(SELECT id, page_id, com_id, title, img, link_type, link, text, start_time, end_time FROM page_component_data_draft WHERE page_id = ?)
		`, req.ID, req.PlatformID).Error; err != nil {
			return err
		}

		tx.Model(&Pages{}).Where("id = ?", req.ID).Update("released_at", gorm.Expr("UNIX_TIMESTAMP()"))
		return nil
	})
	return err
}

/* func (pages *Pages) Validate(platformID int, ctx gin.Context) {
	// data is not exist
	if pages.ID == 0 {
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
} */
