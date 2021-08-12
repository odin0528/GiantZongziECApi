package frontend

import (
	. "eCommerce/internal/database"

	"gorm.io/gorm"
)

type CategoryQuery struct {
	ID         int `json:"id"`
	CustomerID int
	ParentID   int `uri:"parent_id"`
	Sort       int
	DeletedAt  int
}

type Category struct {
	ID         int    `json:"id"`
	CustomerID int    `json:"-"`
	ParentID   int    `json:"parent_id"`
	Layer      int    `json:"layer"`
	Title      string `json:"title"`
	Sort       int    `json:"sort"`
	DeletedAt  int    `json:"-"`
	TimeDefault
}

func (Category) TableName() string {
	return "category"
}

func (query *CategoryQuery) Query() *gorm.DB {
	sql := DB.Table("category").Where("deleted_at = 0")
	if query.ID != 0 {
		sql.Where("id = ?", query.ID)
	}

	if query.ParentID != 0 {
		sql.Where("parent_id = ?", query.ParentID)
	}

	if query.CustomerID != 0 {
		sql.Where("customer_id = ?", query.CustomerID)
	}

	if query.Sort != 0 {
		sql.Where("sort = ?", query.Sort)
	}
	return sql
}

func (query *CategoryQuery) Fetch() (category Category) {
	sql := query.Query()
	sql.First(&category)
	return
}

func (query *CategoryQuery) FetchAll() (categories []Category) {
	sql := query.Query()
	sql.Order("sort asc").Scan(&categories)
	return
}

func (query *CategoryQuery) Count() (count int64) {
	sql := query.Query()
	sql.Count(&count)
	return
}

func (query *CategoryQuery) GetBreadcrumbs(breadcrumbs *[]Category) {
	if query.ParentID != -1 {
		parentCategoryQuery := CategoryQuery{
			ID:         query.ParentID,
			CustomerID: query.CustomerID,
		}
		parentCategory := parentCategoryQuery.Fetch()
		if parentCategory.ID != 0 {
			parentCategoryQuery.ParentID = parentCategory.ParentID
			*breadcrumbs = append([]Category{parentCategory}, *breadcrumbs...)
			parentCategoryQuery.GetBreadcrumbs(breadcrumbs)
		}
	}

}
