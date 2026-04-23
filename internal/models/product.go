package models

// Product representa los datos de un vape o producto del catálogo.
type Product struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`        // Nombre comercial
	Subtitle    string   `json:"subtitle"`    // Descripción corta o lema
	Description string   `json:"description"` // Descripción detallada
	Price       float64  `json:"price"`       // Precio en euros
	Image       string   `json:"image"`       // Ruta a la imagen del producto
	Alt         string   `json:"alt"`         // Texto alternativo para accesibilidad
	Flavors     []string `json:"flavors"`     // Lista de sabores/variantes disponibles
	Featured    bool     `json:"featured"`    // Indica si aparece en la portada
}
