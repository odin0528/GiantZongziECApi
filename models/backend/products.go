package backend

import (
	. "eCommerce/internal/database"

	"github.com/liudng/godump"
	"gorm.io/gorm"
)

type ProductQuery struct {
	ID         int `json:"id" uri:"id"`
	PlatformID int `json:"-"`
}

type ProductListReq struct {
	PlatformID int `json:"-"`
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
	IsPublic        bool                  `json:"is_public"`
	DeletedAt       int                   `json:"-"`
	TimeDefault
}

// 基本CURD功能
func (products *Products) Create() (err error) {
	err = DB.Create(&products).Error
	return
}

func (products *Products) Update() (err error) {
	err = DB.Debug().Save(&products).Error
	return
}

// 資料驗證
func (products *Products) Validate(platformID int) bool {
	// data is not exist or The owner of the data is not the operator
	if products.ID == 0 || products.PlatformID != platformID {
		return false
	}
	return true
}

// 查詢功能
func (query *ProductQuery) Query() *gorm.DB {
	sql := DB.Table("products").Where("deleted_at = 0")
	if query.ID != 0 {
		sql.Where("id = ?", query.ID)
	}

	if query.PlatformID != 0 {
		sql.Where("platform_id = ?", query.PlatformID)
	}

	return sql
}

func (query *ProductQuery) Fetch() (product Products) {
	sql := query.Query()
	sql.First(&product)
	return
}

func (query *ProductQuery) FetchAll() (products []Products) {
	sql := query.Query()
	sql.Scan(&products)
	return
}

func (req *ProductListReq) FetchAll() (products []Products, pagination Pagination) {
	var count int64
	sql := DB.Debug().Table("products").Where("platform_id = ?", req.PlatformID)
	sql.Count(&count)
	sql.Offset((req.Page - 1) * req.Items).Limit(req.Items).Scan(&products)
	pagination = CreatePagination(req.Page, req.Items, count)
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

func (product *Products) ChangePubliced() {
	godump.Dump(product)
	DB.Debug().Table("products").Where("id = ?", product.ID).Update("is_public", product.IsPublic)
}
