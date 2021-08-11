package frontend

type ProductStyle struct {
	ID         int    `json:"id" gorm:"<-:create"`
	CustomerID int    `json:"-" gorm:"<-:create"`
	ProductID  int    `json:"-" gorm:"<-:create"`
	Title      string `json:"title"`
	Img        string `json:"photo"`
	Sort       int    `json:"sort"`
	TimeDefault
}

func (ProductStyle) TableName() string {
	return "product_style"
}
