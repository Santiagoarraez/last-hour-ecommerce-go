package handlers

import (
	"net/http"
	"strings"

	"lasthour/internal/models"
)

// Products muestra el catálogo completo de productos.
func (a *App) Products(w http.ResponseWriter, r *http.Request) {
	products, err := a.products.ListProducts()
	if err != nil {
		http.Error(w, "No se pudieron cargar los productos", http.StatusInternalServerError)
		return
	}

	a.render(w, r, "products.html", map[string]any{
		"Title":    "Nuestro Catálogo - Last Hour",
		"Products": products,
	})
}

// ProductDetail muestra la información detallada de un solo vape.
func (a *App) ProductDetail(w http.ResponseWriter, r *http.Request) {
	// Extraemos el identificador quitando el prefijo de la ruta
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

	// Buscamos otros productos de la misma categoría (sabores del mismo modelo)
	allProducts, _ := a.products.ListProducts()
	var variants []models.Product
	for _, p := range allProducts {
		if p.Category == product.Category && p.ID != product.ID {
			variants = append(variants, p)
		}
	}

	a.render(w, r, "product_detail.html", map[string]any{
		"Title":    product.Name + " - Last Hour",
		"Product":  product,
		"Variants": variants,
	})
}
