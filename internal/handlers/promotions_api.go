package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"lasthour/internal/models"
)

// ApiPromotions actúa como router para /api/promotions (sin ID).
func (a *App) ApiPromotions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.ApiListPromotions(w, r)
	case http.MethodPost:
		a.ApiCreatePromotion(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "método no permitido")
	}
}

// ApiPromotionByID actúa como router para /api/promotions/{id}.
func (a *App) ApiPromotionByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/promotions/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id requerido")
		return
	}

	switch r.Method {
	case http.MethodPut:
		a.ApiUpdatePromotion(w, r, id)
	case http.MethodDelete:
		a.ApiDeletePromotion(w, r, id)
	default:
		writeError(w, http.StatusMethodNotAllowed, "método no permitido")
	}
}

// ApiListPromotions devuelve todas las promociones en JSON.
func (a *App) ApiListPromotions(w http.ResponseWriter, r *http.Request) {
	list, err := a.promotions.ListPromotions()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error al obtener promociones")
		return
	}
	if list == nil {
		list = []models.Promotion{}
	}
	writeJSON(w, http.StatusOK, list)
}

// ApiCreatePromotion crea un nuevo pack promocional.
func (a *App) ApiCreatePromotion(w http.ResponseWriter, r *http.Request) {
	if _, ok := a.requireSellerAPI(w, r); !ok {
		return
	}

	var input struct {
		Name        string                 `json:"name"`
		Description string                 `json:"description"`
		Price       float64                `json:"price"`
		Units       int                    `json:"units"`
		Image       string                 `json:"image"`
		Items       []models.PromotionItem `json:"items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	// Convertimos a string para mantener compatibilidad con el servicio actual
	priceStr := fmt.Sprintf("%.2f", input.Price)
	unitsStr := fmt.Sprintf("%d", input.Units)

	err := a.promotions.CreatePromotion(input.Name, input.Description, priceStr, unitsStr, input.Image, input.Items)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	id := strings.ToLower(strings.ReplaceAll(input.Name, " ", "-"))
	created, _ := a.promotions.FindPromotionByID(id)
	writeJSON(w, http.StatusCreated, created)
}

// ApiUpdatePromotion actualiza una promoción existente.
func (a *App) ApiUpdatePromotion(w http.ResponseWriter, r *http.Request, id string) {
	if _, ok := a.requireSellerAPI(w, r); !ok {
		return
	}

	var input models.Promotion
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}
	input.ID = id // Asegurar que el ID sea el de la URL

	if err := a.promotions.UpdatePromotion(input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	updated, _ := a.promotions.FindPromotionByID(id)
	writeJSON(w, http.StatusOK, updated)
}

// ApiDeletePromotion elimina una promoción por su ID.
func (a *App) ApiDeletePromotion(w http.ResponseWriter, r *http.Request, id string) {
	if _, ok := a.requireSellerAPI(w, r); !ok {
		return
	}

	if err := a.promotions.DeletePromotion(id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
