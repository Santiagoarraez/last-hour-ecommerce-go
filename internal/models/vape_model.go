package models

// VapeModel representa la estructura base de un modelo de vape.
type VapeModel struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Subtitle    string  `json:"subtitle"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Puffs       int     `json:"puffs,omitempty"`
	Battery     string  `json:"battery,omitempty"`
	Liquid      string  `json:"liquid,omitempty"`
	Nicotine    string  `json:"nicotine,omitempty"`
	BestSeller  bool    `json:"best_seller,omitempty"`
}
