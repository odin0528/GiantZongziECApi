package backend

import (
	. "eCommerce/internal/database"
)

type PageComponentDataDraftQuery struct {
	ComID int `json:"com_id"`
}

type PageComponentDataDraft struct {
	ID        int    `json:"-"`
	ComID     int    `json:"-"`
	PageID    int    `json:"-"`
	Title     string `json:"title"`
	Img       string `json:"img"`
	Link      string `json:"link"`
	Text      string `json:"text"`
	StartTime int    `json:"start_time"`
	EndTime   int    `json:"end_time"`
}

func (query *PageComponentDataDraftQuery) FetchByComID() (componentData []PageComponentDataDraft) {
	DB.Model(&PageComponentDataDraft{}).Where("com_id = ?", query.ComID).Scan(&componentData)
	return
}

func (listData *PageComponentDataDraft) Save() (err error) {
	listData.ID = 0
	return DB.Create(&listData).Error
}
