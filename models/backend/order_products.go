package backend

type OrderProducts struct {
	ID              int     `json:"-"`
	OrderID         int     `json:"-"`
	ProductID       int     `json:"-"`
	StyleID         int     `json:"-"`
	Qty             int     `json:"qty"`
	StockQty        int     `json:"stock_qty"`
	Price           float64 `json:"price"`
	DiscountedPrice float64 `json:"discounted_price"`
	Total           float64 `json:"total"`
	Title           string  `json:"title"`
	StyleTitle      string  `json:"style_title"`
	Photo           string  `json:"photo"`
	Sku             string  `json:"sku"`
	TimeDefault
}
