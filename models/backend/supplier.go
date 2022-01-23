package backend

import (
	. "eCommerce/internal/database"
	"eCommerce/internal/rdb"
	"fmt"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

type Supplier struct {
	ID           int    `json:"id" gorm:"<-:create"`
	PlatformID   int    `json:"-" gorm:"<-:create"`
	Title        string `json:"title"`
	Url          string `json:"url"`
	ContactName  string `json:"contact_name"`
	ContactPhone string `json:"contact_phone"`
	Memo         string `json:"memo"`
	DeletedAt    soft_delete.DeletedAt
	TimeDefault
}

type SupplierQuery struct {
	PlatformID int
}

func (Supplier) TableName() string {
	return "supplier"
}

func (query *SupplierQuery) GetCondition() *gorm.DB {
	sql := DB.Model(Supplier{})

	sql.Where("platform_id = ?", query.PlatformID)

	return sql
}

func (query *SupplierQuery) FetchAll() (suppliers []Supplier, err error) {

	key := fmt.Sprintf("platform_%d_supplier", query.PlatformID)
	err = rdb.Get(key, &suppliers)
	if err == redis.Nil {
		sql := query.GetCondition()
		err = sql.Scan(&suppliers).Error

		if err != nil {
			return
		}

		rdb.Set(key, suppliers)
	}
	return
}
