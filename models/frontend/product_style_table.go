package frontend

type ProductStyleTable struct {
	ID         int    `json:"id" gorm:"<-:create"`
	CustomerID int    `json:"-" gorm:"<-:create"`
	ProductID  int    `json:"-" gorm:"<-:create"`
	Group      int    `json:"-"`
	Title      string `json:"title"`
	SubTitle   string `json:"subTitle"`
	Sku        string `json:"sku"`
	Price      int    `json:"price"`
	Qty        int    `json:"qty"`
	DeletedAt  int    `json:"-"`
	TimeDefault
}

func (ProductStyleTable) TableName() string {
	return "product_style_table"
}

// 基本CURD功能
