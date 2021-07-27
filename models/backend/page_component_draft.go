package backend

import (
	. "eCommerce/internal/database"

	"github.com/jinzhu/gorm"
)

type PageComponentDraftReq struct {
	PageID int    `json:"page_id" uri:"id"`
	Type   string `json:"type" uri:"type"`
}

type PageComponentDraft struct {
	ID            int                      `json:"id" gorm:"primaryKey;autoIncrement"`
	CustomerID    int                      `json:"customer_id"`
	PageID        int                      `json:"page_id" uri:"page_id"`
	Sort          int                      `json:"sort"`
	ComponentName string                   `json:"componentName"`
	Title         string                   `json:"title"`
	Type          string                   `json:"type"`
	Data          []PageComponentDataDraft `json:"data" gorm:"-"`
	TimeDefault
}

type PageComponentDraftChangeReq struct {
	PageID    int `json:"page_id"`
	Position1 int `json:"position_1"`
	Position2 int `json:"position_2"`
}

func (req *PageComponentDraftReq) DmoComponentFetch() (component PageComponentDraft) {
	DB.Debug().Model(&component).Where("page_id = ? and type = ?", req.PageID, req.Type).Scan(&component)
	return
}

func (component *PageComponentDraft) FetchByPageID() (components []PageComponentDraft) {
	DB.Debug().Model(&component).Where("page_id = ?", component.PageID).Order("sort asc").Scan(&components)
	return
}

func (component *PageComponentDraft) Save() (err error) {
	component.ID = 0
	err = DB.Table("page_component_draft").Where("page_id = ? AND sort >= ?", component.PageID, component.Sort).Update("sort", gorm.Expr("sort + 1")).Error
	err = DB.Create(&component).Error
	return
}

func (component *PageComponentDraft) FetchBySort() {
	DB.Model(component).Where("page_id = ? AND sort = ?", component.PageID, component.Sort).Scan(&component)
}

func (component *PageComponentDraft) Delete() (err error) {
	DB.Delete(PageComponentDraft{}, "id = ?", component.ID)
	err = DB.Debug().Table("page_component_draft").Where("page_id = ? AND sort > ?", component.PageID, component.Sort).Update("sort", gorm.Expr("sort - 1")).Error
	return
}

func (component *PageComponentDraft) DeleteChildren() (err error) {
	DB.Delete(PageComponentDataDraft{}, "com_id = ?", component.ID)
	return
}

func (component *PageComponentDraft) Update() (err error) {
	err = DB.Save(&component).Error
	return
}
