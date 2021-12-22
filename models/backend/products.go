package backend

import (
	. "eCommerce/internal/database"

	"github.com/liudng/godump"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type ProductQuery struct {
	ID         int    `json:"id" uri:"id"`
	PlatformID int    `json:"-"`
	Title      string `json:"title"`
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
	Min             float64               `json:"min"`
	Max             float64               `json:"max"`
	IsPublic        bool                  `json:"is_public"`
	DeletedAt       soft_delete.DeletedAt `json:"-"`
	TimeDefault
}

// 基本CURD功能
func (products *Products) Create() (err error) {
	err = DB.Create(&products).Error
	return
}

func (products *Products) Update() (err error) {
	err = DB.Save(&products).Error
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
	sql := DB.Model(&Products{})
	if query.ID != 0 {
		sql.Where("id = ?", query.ID)
	}

	if query.Title != "" {
		sql.Where("title like ?", "%"+query.Title+"%")
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

func (query *ProductQuery) FetchAll() (products []Products, pagination Pagination) {
	var count int64
	sql := query.Query()
	sql.Count(&count)
	sql.Offset((query.Page - 1) * query.Items).Limit(query.Items).Order("created_at DESC").Scan(&products)
	pagination = CreatePagination(query.Page, query.Items, count)
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
	DB.Table("product_style_table").Where("product_id = ?", product.ID).Order("group_no ASC, sort ASC").Scan(&styleList)

	for _, style := range styleList {
		if len(product.StyleTable) < style.GroupNo+1 {
			product.StyleTable = append(product.StyleTable, []ProductStyleTable{})
		}
		product.StyleTable[style.GroupNo] = append(product.StyleTable[style.GroupNo], style)
	}
}

func (product *Products) GetLowStockStyleTable() {
	var styleList []ProductStyleTable
	DB.Table("product_style_table").Where("product_id = ? and qty <= low_stock", product.ID).Order("group_no ASC, sort ASC").Scan(&styleList)

	index := -1
	groupNo := -1
	for _, style := range styleList {
		if groupNo != style.GroupNo {
			product.StyleTable = append(product.StyleTable, []ProductStyleTable{})
			groupNo = style.GroupNo
			index++
		}
		product.StyleTable[index] = append(product.StyleTable[index], style)
	}
}

func (product *Products) GetOverSaleStyleTable() {
	var styleList []ProductStyleTable
	DB.Table("product_style_table").Where("product_id = ? and qty < 0", product.ID).Order("group_no ASC, sort ASC").Scan(&styleList)

	index := -1
	groupNo := -1
	for _, style := range styleList {
		if groupNo != style.GroupNo {
			product.StyleTable = append(product.StyleTable, []ProductStyleTable{})
			groupNo = style.GroupNo
			index++
		}
		product.StyleTable[index] = append(product.StyleTable[index], style)
	}
}

func (product *Products) GetWaitDeliveryStyleTable() {
	var styleList []ProductStyleTable
	DB.Table("report_wait_delivery_style_table").Where("product_id = ?", product.ID).Order("group_no ASC, sort ASC").Scan(&styleList)

	index := -1
	groupNo := -1
	for _, style := range styleList {
		if groupNo != style.GroupNo {
			product.StyleTable = append(product.StyleTable, []ProductStyleTable{})
			groupNo = style.GroupNo
			index++
		}
		product.StyleTable[index] = append(product.StyleTable[index], style)
	}
}

func (product *Products) ChangePubliced() {
	godump.Dump(product)
	DB.Table("products").Where("id = ?", product.ID).Update("is_public", product.IsPublic)
}
