package backend

type TimeDefault struct {
	CreatedAt int `json:"created_at" gorm:"<-:create"`
	UpdatedAt int `json:"updated_at"`
}

type Pagination struct {
	Page  int `json:"page"`
	Start int `json:"start"`
	End   int `json:"end"`
	Total int `json:"total"`
	Items int `json:"items"`
}

func CreatePagination(page int, items int, total int64) (pagination Pagination) {
	pagination.Page = page
	pagination.Start = ((page - 1) * items) + 1
	pagination.End = pagination.Start + items
	pagination.Total = int(total)
	pagination.Items = items
	if pagination.End > pagination.Total {
		pagination.End = pagination.Total
	}
	return
}
