package backend

import (
	. "eCommerce/internal/database"

	"gorm.io/gorm"
)

// DB Table struct
// 不是新增的情況都不能直接bingJSON
type PageComponentDraft struct {
	ID            int                      `json:"-" gorm:"primaryKey;autoIncrement"`
	PlatformID    int                      `json:"-"`
	PageID        int                      `json:"page_id"`
	Sort          int                      `json:"sort"`
	ComponentName string                   `json:"componentName"`
	Title         string                   `json:"title"`
	Text          string                   `json:"text"`
	Type          string                   `json:"type"`
	Data          []PageComponentDataDraft `json:"data" gorm:"-"`
	TimeDefault   `json:"-"`
}

func (PageComponentDraft) TableName() string {
	return "page_component_draft"
}

// 基本CURD功能
func (component *PageComponentDraft) Create() (err error) {
	component.ID = 0
	err = DB.Table("page_component_draft").Where("page_id = ? AND platform_id = ? AND sort >= ?", component.PageID, component.PlatformID, component.Sort).Update("sort", gorm.Expr("sort + 1")).Error
	err = DB.Create(&component).Error
	return
}

func (component *PageComponentDraft) Update() (err error) {
	err = DB.Save(&component).Error
	return
}

func (component *PageComponentDraft) Delete() (err error) {
	DB.Delete(PageComponentDraft{}, "id = ?", component.ID)
	err = DB.Table("page_component_draft").Where("page_id = ? AND platform_id = ? AND sort > ?", component.PageID, component.PlatformID, component.Sort).Update("sort", gorm.Expr("sort - 1")).Error
	return
}

func (component *PageComponentDraft) DeleteChildren() (err error) {
	DB.Delete(PageComponentDataDraft{}, "com_id = ?", component.ID)
	return
}

func (component *PageComponentDraft) Validate(platformID int) bool {
	// data is not exist or The owner of the data is not the operator
	if component.ID == 0 || component.PlatformID != platformID {
		return false
	}
	return true
}

//一般查詢功能
type PageComponentDraftQuery struct {
	PageID     int `json:"page_id"`
	PlatformID int `json:"platform_id"`
	Sort       int `json:"sort"`
}

func (req *PageComponentDraftQuery) Fetch() (component PageComponentDraft) {
	DB.Table("page_component_draft").Where("page_id = ? and sort = ? AND platform_id = ?", req.PageID, req.Sort, req.PlatformID).Scan(&component)
	return
}

func (req *PageComponentDraftQuery) FetchByPageID() (components []PageComponentDraft) {
	DB.Table("page_component_draft").Select("id, sort, component_name, title, text, type").
		Where("page_id = ? AND platform_id = ?", req.PageID, req.PlatformID).Order("sort asc").Scan(&components)
	return
}

// 位置交換請求
type PageComponentDraftChangeReq struct {
	PageID  int `json:"page_id"`
	Sort    int `json:"sort"`
	NewSort int `json:"new_sort"`
}

// 編輯功能請求
type PageComponentDraftEditReq struct {
	PageID int                `json:"page_id"`
	Sort   int                `json:"sort"`
	Data   PageComponentDraft `json:"data"`
}
