package frontend

import (
	. "eCommerce/internal/database"

	"gorm.io/gorm"
)

type ProductQuery struct {
	ID         int    `json:"id" uri:"id"`
	CategoryID int    `json:"category_id" uri:"category_id"`
	Layer      int    `json:"layer" uri:"layer"`
	PlatformID int    `json:"-"`
	Min        int    `json:"min"`
	Max        int    `json:"max"`
	Sort       string `json:"sort"`
	Pagination
}

type Products struct {
	ID              int                   `json:"id"`
	PlatformID      int                   `json:"-" gorm:"<-:create"`
	Title           string                `json:"title"`
	StyleTitle      string                `json:"style_title"`
	SubStyleTitle   string                `json:"sub_style_title"`
	Description     string                `json:"description"`
	CategoryLayer1  int                   `json:"category_layer1"`
	CategoryLayer2  int                   `json:"category_layer2"`
	CategoryLayer3  int                   `json:"category_layer3"`
	CategoryLayer4  int                   `json:"category_layer4"`
	Photos          []ProductPhotos       `json:"photos" gorm:"-"`
	Style           []ProductStyle        `json:"style" gorm:"-"`
	SubStyle        []ProductSubStyle     `json:"sub_style" gorm:"-"`
	StyleTable      [][]ProductStyleTable `json:"style_table" gorm:"-"`
	StyleEnabled    bool                  `json:"style_enabled"`
	SubStyleEnabled bool                  `json:"sub_style_enabled"`
	MinPrice        int                   `json:"min"`
	MaxPrice        int                   `json:"max"`
	IsPublic        bool                  `json:"is_public"`
	DeletedAt       int                   `json:"-"`
	TimeDefault
}

// 查詢功能
func (query *ProductQuery) Query() *gorm.DB {
	sql := DB.Table("products").Where("deleted_at = 0 AND is_public = 1")
	if query.ID != 0 {
		sql.Where("products.id = ?", query.ID)
	}

	if query.CategoryID > 0 && query.Layer > 0 {
		switch query.Layer {
		case 1:
			sql.Where("category_layer1 = ?", query.CategoryID)
		case 2:
			sql.Where("category_layer2 = ?", query.CategoryID)
		case 3:
			sql.Where("category_layer3 = ?", query.CategoryID)
		case 4:
			sql.Where("category_layer4 = ?", query.CategoryID)
		}
	}

	if query.PlatformID != 0 {
		sql.Where("products.platform_id = ?", query.PlatformID)
	}

	if query.Min > 0 || query.Max > 0 || query.Sort == "price-asc" || query.Sort == "price-desc" {
		sql.Joins("inner join product_style_table on products.id = product_style_table.product_id")
		if query.Min > 0 {
			sql.Where("product_style_table.price >= ?", query.Min)
		}

		if query.Max > 0 {
			sql.Where("product_style_table.price <= ?", query.Max)
		}

		if query.Sort == "price-asc" {
			sql.Order("product_style_table.price ASC")
		} else if query.Sort == "price-desc" {
			sql.Order("product_style_table.price DESC")
		}

		sql.Group("products.id")
	}

	return sql
}

func (query *ProductQuery) Fetch() (product Products) {
	sql := query.Query()
	sql.First(&product)
	return
}

func (query *ProductQuery) FetchAll() (products []Products, pagination Pagination) {
	var count int64
	sql := query.Query()
	sql.Count(&count)
	sql.Select("products.*").Order("products.created_at DESC").Offset((query.Page - 1) * Items).Limit(Items).Scan(&products)
	pagination = CreatePagination(query.Page, Items, count)
	return
}

func (query *ProductQuery) FetchRelated() (products []Products) {
	id := query.ID
	query.ID = 0
	sql := query.Query()
	sql.Where("products.id != ?", id).Limit(10).Scan(&products)
	return
}

// 關連功能
func (product *Products) GetPhotos() {
	DB.Table("product_photos").Where("product_id = ?", product.ID).Order("sort ASC").Scan(&product.Photos)
}

func (product *Products) GetStyle() {
	DB.Table("product_style").Where("product_id = ?", product.ID).Order("sort ASC").Scan(&product.Style)
}

func (product *Products) GetSubStyle() {
	DB.Table("product_sub_style").Where("product_id = ?", product.ID).Order("sort ASC").Scan(&product.SubStyle)
}

func (product *Products) GetStyleTable() {
	var styleList []ProductStyleTable
	DB.Table("product_style_table").Where("product_id = ?", product.ID).Scan(&styleList)

	for _, style := range styleList {
		if len(product.StyleTable) < style.Group+1 {
			product.StyleTable = append(product.StyleTable, []ProductStyleTable{})
		}
		product.StyleTable[style.Group] = append(product.StyleTable[style.Group], style)
	}
}

func (product *Products) GetRelated() (products []Products) {
	query := ProductQuery{
		ID:         product.ID,
		PlatformID: product.PlatformID,
	}

	if product.CategoryLayer4 != -1 {
		query.Layer = 4
		query.CategoryID = product.CategoryLayer4
	} else if product.CategoryLayer3 != -1 {
		query.Layer = 3
		query.CategoryID = product.CategoryLayer3
	} else if product.CategoryLayer2 != -1 {
		query.Layer = 2
		query.CategoryID = product.CategoryLayer2
	} else if product.CategoryLayer1 != -1 {
		query.Layer = 1
		query.CategoryID = product.CategoryLayer1
	}

	return query.FetchRelated()
}

func (product *Products) GetPriceRange() {
	DB.Table("product_style_table").Select([]string{"max(price) as max_price", "min(price) as min_price"}).Where("product_id = ?", product.ID).Scan(&product)
}
