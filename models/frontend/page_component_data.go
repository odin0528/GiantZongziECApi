package frontend

import (
	. "eCommerce/internal/database"
)

type PageComponentDataQuery struct {
	ComID int `json:"com_id"`
}

type PageComponentData struct {
	ID        int    `json:"-"`
	ComID     int    `json:"-"`
	PageID    int    `json:"-"`
	Title     string `json:"title"`
	Img       string `json:"img"`
	Link      string `json:"link"`
	LinkType  int    `json:"link_type"`
	Text      string `json:"text"`
	StartTime int    `json:"start_time"`
	EndTime   int    `json:"end_time"`
	TimeDefault
}

func (PageComponentData) TableName() string {
	return "page_component_data"
}

func (query *PageComponentDataQuery) FetchByComID() (componentData []PageComponentData) {
	DB.Model(&PageComponentData{}).Where("com_id = ?", query.ComID).Scan(&componentData)
	return
}
