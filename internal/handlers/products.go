package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"

	"lasthour/internal/models"
)

// Products muestra el catálogo completo organizado por Modelos, Sabores y Promociones.
func (a *App) Products(w http.ResponseWriter, r *http.Request) {
	// Cargamos los datos desde los nuevos servicios modulares
	modelsList, _ := a.vapeModels.ListModels()
	flavorsList, _ := a.flavors.ListFlavors()
	promosList, _ := a.promotions.ListPromotions()

	// Convertimos a JSON para que el frontend pueda usarlos en el modal de promos
	modelsJson, _ := json.Marshal(modelsList)
	flavorsJson, _ := json.Marshal(flavorsList)
	promosJson, _ := json.Marshal(promosList)

	a.render(w, r, "products.html", map[string]any{
		"Title":           "Nuestro Catálogo - Last Hour",
		"Models":          modelsList,
		"Flavors":         flavorsList,
		"Promotions":      promosList,
		"ModelsJSON":      template.JS(modelsJson),
		"FlavorsJSON":     template.JS(flavorsJson),
		"PromotionsJSON":  template.JS(promosJson),
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
