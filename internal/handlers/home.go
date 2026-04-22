package handlers

import (
	"net/http"

	"lasthour/internal/models"
)

type HomePageData struct {
	Title    string
	Products []models.Product
}

func (a *App) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	products, err := a.products.ListFeaturedProducts()
	if err != nil {
		http.Error(w, "No se pudieron cargar los productos", http.StatusInternalServerError)
		return
	}

	a.render(w, "home.html", HomePageData{
		Title:    "Vape Store",
		Products: products,
	})
}
