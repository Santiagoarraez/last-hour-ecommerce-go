package handlers

import (
	"net/http"
	"strings"

	"lasthour/internal/models"
)

// FlavorDetail maneja la visualización detallada de un sabor específico.
// GET /products/flavor/{id}
func (a *App) FlavorDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extraemos el ID del sabor de la URL
	flavorID := strings.TrimPrefix(r.URL.Path, "/products/flavor/")
	if flavorID == "" {
		http.NotFound(w, r)
		return
	}

	// 1. Buscamos el sabor
	flavor, err := a.flavors.FindFlavorByID(flavorID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// 2. Buscamos el modelo al que pertenece
	model, err := a.vapeModels.FindModelByID(flavor.ModelID)
	if err != nil {
		http.Error(w, "Modelo no encontrado", http.StatusInternalServerError)
		return
	}

	// 3. Buscamos los otros sabores del mismo modelo (Siblings) para los thumbnails
	allFlavors, _ := a.flavors.ListFlavors()
	var siblings []models.Flavor
	for _, f := range allFlavors {
		if f.ModelID == flavor.ModelID && f.ID != flavor.ID {
			siblings = append(siblings, f)
		}
	}

	// 4. Renderizamos la vista de detalle
	a.render(w, r, "product_detail.html", map[string]any{
		"Title":    model.Name + " " + flavor.Name + " - Last Hour",
		"Flavor":   flavor,
		"Model":    model,
		"Siblings": siblings,
	})
}
