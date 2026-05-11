package models

// Flavor representa una variante de sabor para un modelo de vape específico.
type Flavor struct {
	ID          string `json:"id"`
	ModelID     string `json:"model_id"`
	ModelName   string `json:"model_name"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	OutOfStock  bool   `json:"out_of_stock"`
}
