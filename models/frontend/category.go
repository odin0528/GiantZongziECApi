package frontend

import (
	. "eCommerce/internal/database"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type CategoryQuery struct {
	ID         int
	PlatformID int
	ParentID   int `uri:"parent_id"`
	Sort       int
	DeletedAt  int
}

type Category struct {
	ID         int    `json:"id"`
	PlatformID int    `json:"-"`
	ParentID   int    `json:"parent_id"`
	Layer      int    `json:"layer"`
	Title      string `json:"title"`
	Sort       int    `json:"sort"`
	DeletedAt  soft_delete.DeletedAt
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

	if query.PlatformID != 0 {
		sql.Where("platform_id = ?", query.PlatformID)
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
	DB.Debug().Model(&Category{}).
		Where("(parent_id = ? OR parent_id = ?) AND platform_id = ?", query.ID, query.ParentID, query.PlatformID).
		Order(fmt.Sprintf("id = %d DESC, parent_id = %d, sort asc", query.ParentID, query.ID)).Scan(&categories)
	return
}

func (query *CategoryQuery) GetBreadcrumbs(breadcrumbs *[]Category) {
	if query.ParentID != -1 {
		parentCategoryQuery := CategoryQuery{
			ID:         query.ParentID,
			PlatformID: query.PlatformID,
		}
		parentCategory := parentCategoryQuery.Fetch()
		if parentCategory.ID != 0 {
			parentCategoryQuery.ParentID = parentCategory.ParentID
			*breadcrumbs = append([]Category{parentCategory}, *breadcrumbs...)
			parentCategoryQuery.GetBreadcrumbs(breadcrumbs)
		}
	}

}
