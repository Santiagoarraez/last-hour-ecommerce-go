package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"lasthour/internal/models"
)

// ApiModels actúa como router para /api/models (sin ID).
func (a *App) ApiModels(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.ApiListModels(w, r)
	case http.MethodPost:
		a.ApiCreateModel(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "método no permitido")
	}
}

// ApiModelByID actúa como router para /api/models/{id}.
func (a *App) ApiModelByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/models/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id requerido")
		return
	}

	switch r.Method {
	case http.MethodPut:
		a.ApiUpdateModel(w, r, id)
	case http.MethodDelete:
		a.ApiDeleteModel(w, r, id)
	default:
		writeError(w, http.StatusMethodNotAllowed, "método no permitido")
	}
}

// ApiListModels devuelve todos los modelos en JSON.
func (a *App) ApiListModels(w http.ResponseWriter, r *http.Request) {
	list, err := a.vapeModels.ListModels()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error al obtener modelos")
		return
	}
	if list == nil {
		list = []models.VapeModel{}
	}
	writeJSON(w, http.StatusOK, list)
}

// ApiCreateModel crea un nuevo modelo de vape.
func (a *App) ApiCreateModel(w http.ResponseWriter, r *http.Request) {
	if _, ok := a.requireSellerAPI(w, r); !ok {
		return
	}

	var input struct {
		Name        string `json:"name"`
		Subtitle    string `json:"subtitle"`
		Description string `json:"description"`
		Price       float64 `json:"price"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	err := a.vapeModels.CreateModel(input.Name, input.Subtitle, input.Description, input.Price)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Devolver el modelo recién creado (lo buscamos por el ID generado)
	id := strings.ToLower(strings.ReplaceAll(input.Name, " ", "-"))
	created, err := a.vapeModels.FindModelByID(id)
	if err != nil {
		writeJSON(w, http.StatusCreated, map[string]string{"message": "modelo creado"})
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

// ApiUpdateModel actualiza un modelo existente.
func (a *App) ApiUpdateModel(w http.ResponseWriter, r *http.Request, id string) {
	if _, ok := a.requireSellerAPI(w, r); !ok {
		return
	}

	var input struct {
		Name        string `json:"name"`
		Subtitle    string `json:"subtitle"`
		Description string `json:"description"`
		Price       float64 `json:"price"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	// Usamos UpdateModel del servicio para activar la actualización en cascada de sabores
	if err := a.vapeModels.UpdateModel(id, input.Name, input.Subtitle, input.Description, input.Price); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	updated, _ := a.vapeModels.FindModelByID(id)
	writeJSON(w, http.StatusOK, updated)
}

// ApiDeleteModel elimina un modelo por su ID.
func (a *App) ApiDeleteModel(w http.ResponseWriter, r *http.Request, id string) {
	if _, ok := a.requireSellerAPI(w, r); !ok {
		return
	}

	if err := a.vapeModels.DeleteModel(id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
