package backend

type TimeDefault struct {
	CreatedAt int `json:"created_at" gorm:"<-:create"`
	UpdatedAt int `json:"updated_at"`
}

type Pagination struct {
	CurrentPage  int `json:"current_page" uri:"page"`
	TotalItem    int `json:"total_item"`
	LastPage     int `json:"last_page"`
	ItemsPerPage int `json:"item_per_page"`
}
