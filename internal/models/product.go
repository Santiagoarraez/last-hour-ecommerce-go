package models

type Product struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Subtitle    string   `json:"subtitle"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	Image       string   `json:"image"`
	Alt         string   `json:"alt"`
	Flavors     []string `json:"flavors"`
	Featured    bool     `json:"featured"`
}
