package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"lasthour/internal/models"
)

// ApiFlavors actúa como router para /api/flavors (sin ID).
func (a *App) ApiFlavors(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.ApiListFlavors(w, r)
	case http.MethodPost:
		a.ApiCreateFlavor(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "método no permitido")
	}
}

// ApiFlavorByID actúa como router para /api/flavors/{id}.
func (a *App) ApiFlavorByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/flavors/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id requerido")
		return
	}

	switch r.Method {
	case http.MethodPut:
		a.ApiUpdateFlavor(w, r, id)
	case http.MethodDelete:
		a.ApiDeleteFlavor(w, r, id)
	default:
		writeError(w, http.StatusMethodNotAllowed, "método no permitido")
	}
}

// ApiListFlavorsByModel responde GET /api/models/{id}/flavors.
func (a *App) ApiListFlavorsByModel(w http.ResponseWriter, r *http.Request) {
	// Extraer ID del modelo de la URL /api/models/{id}/flavors
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		writeError(w, http.StatusBadRequest, "ID de modelo requerido")
		return
	}
	modelID := parts[3]

	list, err := a.flavors.ListFlavorsByModel(modelID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error al obtener sabores")
		return
	}
	if list == nil {
		list = []models.Flavor{}
	}
	writeJSON(w, http.StatusOK, list)
}

// ApiListFlavors devuelve todos los sabores en JSON.
func (a *App) ApiListFlavors(w http.ResponseWriter, r *http.Request) {
	list, err := a.flavors.ListFlavors()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error al obtener sabores")
		return
	}
	if list == nil {
		list = []models.Flavor{}
	}
	writeJSON(w, http.StatusOK, list)
}

// ApiCreateFlavor crea un nuevo sabor para un modelo.
func (a *App) ApiCreateFlavor(w http.ResponseWriter, r *http.Request) {
	if _, ok := a.requireSellerAPI(w, r); !ok {
		return
	}

	var input struct {
		ModelID   string `json:"modelID"`
		ModelName string `json:"modelName"`
		Name      string `json:"name"`
		Image     string `json:"image"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	err := a.flavors.CreateFlavor(input.ModelID, input.ModelName, input.Name, input.Image)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Generar ID para buscar el objeto creado
	id := strings.ToLower(input.ModelID + "-" + input.Name)
	id = strings.ReplaceAll(id, " ", "-")
	
	created, _ := a.flavors.FindFlavorByID(id)
	writeJSON(w, http.StatusCreated, created)
}

// ApiDeleteFlavor elimina un sabor por su ID.
func (a *App) ApiDeleteFlavor(w http.ResponseWriter, r *http.Request, id string) {
	if _, ok := a.requireSellerAPI(w, r); !ok {
		return
	}

	if err := a.flavors.DeleteFlavor(id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ApiUpdateFlavor actualiza un sabor existente.
func (a *App) ApiUpdateFlavor(w http.ResponseWriter, r *http.Request, id string) {
	if _, ok := a.requireSellerAPI(w, r); !ok {
		return
	}

	var input struct {
		Name  string `json:"name"`
		Image string `json:"image"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := a.flavors.UpdateFlavor(id, input.Name, input.Image); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	updated, _ := a.flavors.FindFlavorByID(id)
	writeJSON(w, http.StatusOK, updated)
}
