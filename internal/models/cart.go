package models

// CartItem representa un producto específico y su cantidad dentro del almacenamiento JSON.
type CartItem struct {
	ProductID  string   `json:"product_id"`
	FlavorID   string   `json:"flavor_id"`
	FlavorName string   `json:"flavor_name"`
	Price      float64  `json:"price"`
	Image      string   `json:"image"`
	Quantity   int      `json:"quantity"`
	Flavors    []string `json:"flavors"` // Sabores seleccionados (para compatibilidad)
}

// Cart es la estructura raíz que vincula a un usuario con su lista de productos seleccionados.
type Cart struct {
	UserID string     `json:"user_id"`
	Items  []CartItem `json:"items"`
}

// CartViewItem se utiliza para mostrar el carrito al usuario, incluyendo detalles completos del producto.
type CartViewItem struct {
	Product  Product  `json:"product"`
	Quantity int      `json:"quantity"`
	Flavors  []string `json:"flavors"`
	Subtotal float64  `json:"subtotal"`
}

// CartView es la vista final del carrito que se pasa a las plantillas HTML.
type CartView struct {
	Items []CartViewItem `json:"items"`
	Total float64        `json:"total"`
}
