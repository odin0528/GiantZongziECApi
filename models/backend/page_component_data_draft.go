package backend

import (
	. "eCommerce/internal/database"
)

type PageComponentDataDraft struct {
	ID        int    `json:"id"`
	ComID     int    `json:"com_id"`
	PageID    int    `json:"page_id"`
	Title     string `json:"title"`
	Img       string `json:"img"`
	Link      string `json:"link"`
	Text      string `json:"text"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

func (comData *PageComponentDataDraft) FetchByComID() {
	DB.Model(&comData).Where("page_id = ?", comData.PageID).Scan(&comData)
}

func (listData *PageComponentDataDraft) Save() (err error) {
	listData.ID = 0
	return DB.Create(&listData).Error
}
