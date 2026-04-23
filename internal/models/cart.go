package models

// CartItem representa un producto específico y su cantidad dentro del almacenamiento JSON.
type CartItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

// Cart es la estructura raíz que vincula a un usuario con su lista de productos seleccionados.
type Cart struct {
	UserID string     `json:"user_id"`
	Items  []CartItem `json:"items"`
}

// CartViewItem se utiliza para mostrar el carrito al usuario, incluyendo detalles completos del producto.
type CartViewItem struct {
	Product  Product
	Quantity int
	Subtotal float64
}

// CartView es la vista final del carrito que se pasa a las plantillas HTML.
type CartView struct {
	Items []CartViewItem
	Total float64
}
