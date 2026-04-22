package handlers

import (
	"net/http"
	"strings"

	"lasthour/internal/models"
)

type ProductsPageData struct {
	Title    string
	Products []models.Product
}

type ProductDetailPageData struct {
	Title   string
	Product models.Product
}

func (a *App) Products(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/products" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	products, err := a.products.ListProducts()
	if err != nil {
		http.Error(w, "No se pudo cargar el catalogo", http.StatusInternalServerError)
		return
	}

	a.render(w, "products.html", ProductsPageData{
		Title:    "Our Products - Last Hour",
		Products: products,
	})
}

func (a *App) ProductDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/products/")
	if id == "" || strings.Contains(id, "/") {
		http.NotFound(w, r)
		return
	}

	product, err := a.products.FindProductByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	a.render(w, "product_detail.html", ProductDetailPageData{
		Title:   product.Name + " - Last Hour",
		Product: product,
	})
}
