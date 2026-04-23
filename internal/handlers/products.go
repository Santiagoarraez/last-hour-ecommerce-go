package handlers

import (
	"net/http"
	"strings"

	"lasthour/internal/models"
)

// ProductsPageData estructura los datos para la página del catálogo completo.
type ProductsPageData struct {
	Title    string
	Products []models.Product
}

// ProductDetailPageData estructura los datos para la ficha de un producto individual.
type ProductDetailPageData struct {
	Title   string
	Product models.Product
}

// Products maneja la visualización de todo el catálogo de productos disponibles.
func (a *App) Products(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Recuperamos todos los productos del servicio
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

// ProductDetail gestiona la visualización de la página de detalle de un producto específico.
// Extrae el ID de la URL y busca el producto correspondiente.
func (a *App) ProductDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extraemos el identificador quitando el prefijo de la ruta
	id := strings.TrimPrefix(r.URL.Path, "/products/")
	if id == "" || strings.Contains(id, "/") {
		http.NotFound(w, r)
		return
	}

	// Buscamos el producto por su ID
	product, err := a.products.FindProductByID(id)
	if err != nil {
		// Si no existe, devolvemos un 404 Not Found
		http.NotFound(w, r)
		return
	}

	a.render(w, "product_detail.html", ProductDetailPageData{
		Title:   product.Name + " - Last Hour",
		Product: product,
	})
}
