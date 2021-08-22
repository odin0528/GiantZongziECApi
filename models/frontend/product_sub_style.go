package frontend

type ProductSubStyle struct {
	ID         int    `json:"id" gorm:"<-:create"`
	PlatformID int    `json:"-" gorm:"<-:create"`
	ProductID  int    `json:"-" gorm:"<-:create"`
	Title      string `json:"title"`
	Sort       int    `json:"sort"`
	TimeDefault
}

func (ProductSubStyle) TableName() string {
	return "product_sub_style"
}

// 基本CURD功能
