package models

// Promotion representa un pack promocional o bundle que incluye varios productos.
type Promotion struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Price       float64         `json:"price"`
	Image       string          `json:"image"`
	Units       int             `json:"units"`
	Items       []PromotionItem `json:"items"`
}

// PromotionItem representa un elemento individual dentro de una promoción.
type PromotionItem struct {
	ModelID   string `json:"model_id"`
	ModelName string `json:"model_name"`
}
