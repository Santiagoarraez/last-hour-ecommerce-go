package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
)

// Home maneja la petición a la página principal de la tienda.
// Filtra únicamente la raíz "/" y obtiene los datos necesarios para la home.
func (a *App) Home(w http.ResponseWriter, r *http.Request) {
	// Verificamos que sea exactamente la raíz
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	// PEC 3: Cargamos los datos desde los nuevos servicios modulares
	modelsList, _ := a.vapeModels.ListModels()
	flavorsList, _ := a.flavors.ListFlavors()
	promosList, _ := a.promotions.ListPromotions()

	// Convertimos a JSON para que el JS pueda usarlos en el modal de promos
	modelsJson, _ := json.Marshal(modelsList)
	flavorsJson, _ := json.Marshal(flavorsList)
	promosJson, _ := json.Marshal(promosList)

	// Renderizamos la plantilla home.html inyectando los datos
	a.render(w, r, "home.html", map[string]any{
		"Title":          "Vape Store - Last Hour",
		"Models":         modelsList,
		"Flavors":        flavorsList,
		"Promotions":     promosList,
		"ModelsJSON":     template.JS(modelsJson),
		"FlavorsJSON":    template.JS(flavorsJson),
		"PromotionsJSON": template.JS(promosJson),
	})
}
