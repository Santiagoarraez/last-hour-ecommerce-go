package handlers

import (
	"net/http"

	"lasthour/internal/models"
)

// HomePageData define los datos que se pasan a la plantilla de inicio.
type HomePageData struct {
	Title    string
	Products []models.Product
}

// Home maneja la petición a la página principal de la tienda.
// Filtra únicamente la raíz "/" y obtiene los productos destacados para mostrar.
func (a *App) Home(w http.ResponseWriter, r *http.Request) {
	// Verificamos que sea exactamente la raíz para evitar capturar rutas inexistentes
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtenemos del servicio los productos que irán en la home
	products, err := a.products.ListFeaturedProducts()
	if err != nil {
		http.Error(w, "No se pudieron cargar los productos", http.StatusInternalServerError)
		return
	}

	// Renderizamos la plantilla home.html inyectando los productos
	a.render(w, "home.html", HomePageData{
		Title:    "Vape Store",
		Products: products,
	})
}
