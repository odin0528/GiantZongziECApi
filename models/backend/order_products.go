package backend

type OrderProducts struct {
	ID              int     `json:"-"`
	OrderID         int     `json:"-"`
	ProductID       int     `json:"-"`
	StyleID         int     `json:"-"`
	Qty             int     `json:"qty"`
	Price           float32 `json:"price"`
	DiscountedPrice float32 `json:"discounted_price"`
	Total           float32 `json:"total"`
	Title           string  `json:"title"`
	StyleTitle      string  `json:"style_title"`
	Photo           string  `json:"photo"`
	Sku             string  `json:"sku"`
	TimeDefault
}
