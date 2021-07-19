package backend

import "time"

type TimeDefault struct {
	CreatedAt *time.Time `json:"-" gorm:"column:CreatedAt"` //建立時間
	UpdatedAt *time.Time `json:"-" gorm:"column:UpdatedAt"` //修改時間
	DeletedAt *time.Time `json:"-" gorm:"column:DeletedAt"` //停用時間
}

type Pagination struct {
	CurrentPage  int `json:"current_page" uri:"page"`
	TotalItem    int `json:"total_item"`
	LastPage     int `json:"last_page"`
	ItemsPerPage int `json:"item_per_page"`
}
