package frontend

import (
	. "eCommerce/internal/database"
)

// DB Table struct
// 不是新增的情況都不能直接bingJSON
type PageComponent struct {
	ID            int                 `json:"-" gorm:"primaryKey;autoIncrement"`
	CustomerID    int                 `json:"-"`
	PageID        int                 `json:"page_id"`
	Sort          int                 `json:"sort"`
	ComponentName string              `json:"componentName"`
	Title         string              `json:"title"`
	Text          string              `json:"text"`
	Type          string              `json:"type"`
	Data          []PageComponentData `json:"data" gorm:"-"`
	TimeDefault   `json:"-"`
}

func (component *PageComponent) Validate(customerID int) bool {
	// data is not exist or The owner of the data is not the operator
	if component.ID == 0 || component.CustomerID != customerID {
		return false
	}
	return true
}

//一般查詢功能
type PageComponentQuery struct {
	PageID     int `json:"page_id"`
	CustomerID int `json:"customer_id"`
	Sort       int `json:"sort"`
}

func (req *PageComponentQuery) Fetch() (component PageComponent) {
	DB.Table("page_component_draft").Where("page_id = ? and sort = ? AND customer_id = ?", req.PageID, req.Sort, req.CustomerID).Scan(&component)
	return
}

func (req *PageComponentQuery) FetchByPageID() (components []PageComponent) {
	DB.Table("page_component_draft").Select("id, sort, component_name, title, text, type").
		Where("page_id = ? AND customer_id = ?", req.PageID, req.CustomerID).Order("sort asc").Scan(&components)
	return
}

// 位置交換請求
type PageComponentChangeReq struct {
	PageID  int `json:"page_id"`
	Sort    int `json:"sort"`
	NewSort int `json:"new_sort"`
}

// 編輯功能請求
type PageComponentEditReq struct {
	PageID int           `json:"page_id"`
	Sort   int           `json:"sort"`
	Data   PageComponent `json:"data"`
}
