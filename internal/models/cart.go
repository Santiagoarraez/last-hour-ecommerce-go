package models

type CartItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type Cart struct {
	UserID string     `json:"user_id"`
	Items  []CartItem `json:"items"`
}

type CartViewItem struct {
	Product  Product
	Quantity int
	Subtotal float64
}

type CartView struct {
	Items []CartViewItem
	Total float64
}
