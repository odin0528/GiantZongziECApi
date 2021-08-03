package backend

import (
	. "eCommerce/internal/database"
	"time"

	"gorm.io/gorm"
)

type CategoryQuery struct {
	ID         int `json:"id"`
	CustomerID int
	ParentID   int `uri:"parent_id"`
	Sort       int
	DeletedAt  int
}

type CategoryModifyReq struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type CategoryMoveReq struct {
	ParentID  int `uri:"parent_id"`
	Sort      int `json:"sort"`
	Direction int `json:"direction"`
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

// 基本CURD功能
func (category *Category) Create() (err error) {
	category.ID = 0
	err = DB.Create(&category).Error
	return
}

func (category *Category) Update() (err error) {
	err = DB.Debug().Save(&category).Error
	return
}

func (category *Category) Delete() (err error) {
	category.DeletedAt = int(time.Now().Unix())
	category.Update()
	err = DB.Table("category").Where("parent_id = ? AND customer_id = ? AND sort > ?", category.ParentID, category.CustomerID, category.Sort).Update("sort", gorm.Expr("sort - 1")).Error
	return
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
	sql.Debug().First(&category)
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
			ID: query.ParentID,
		}
		parentCategory := parentCategoryQuery.Fetch()
		parentCategoryQuery.ParentID = parentCategory.ParentID
		*breadcrumbs = append(*breadcrumbs, parentCategory)
		parentCategoryQuery.GetBreadcrumbs(breadcrumbs)
	}

}

func (category *Category) Validate(customerID int) bool {
	if category.ID == 0 || category.CustomerID != customerID {
		return false
	}
	return true
}
