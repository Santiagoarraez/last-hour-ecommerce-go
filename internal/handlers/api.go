package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"lasthour/internal/models"
)

// ---- Helpers de la API REST ----

// writeJSON serializa cualquier valor a JSON y lo envía con el código de estado indicado.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError envía un mensaje de error en formato JSON.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// requireSellerAPI es la versión del middleware requireSeller para la API REST.
// En vez de redirigir, devuelve un JSON 401 o 403.
func (a *App) requireSellerAPI(w http.ResponseWriter, r *http.Request) (models.User, bool) {
	user, ok := a.currentUser(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "no autenticado")
		return models.User{}, false
	}
	if !user.IsSeller() {
		writeError(w, http.StatusForbidden, "acceso restringido a vendedores")
		return models.User{}, false
	}
	return user, true
}

// ---- Endpoints de la API REST ----

// ApiProducts actúa como router para /api/products (sin ID).
// GET  → lista todos los productos (público)
// POST → crea un nuevo producto (solo seller)
func (a *App) ApiProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.ApiListProducts(w, r)
	case http.MethodPost:
		a.ApiCreateProduct(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "método no permitido")
	}
}

// ApiProductByID actúa como router para /api/products/{id}.
// GET    → devuelve un producto por ID (público)
// PUT    → actualiza un producto (solo seller)
// DELETE → elimina un producto (solo seller)
func (a *App) ApiProductByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/products/")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id requerido")
		return
	}

	switch r.Method {
	case http.MethodGet:
		a.ApiGetProduct(w, r, id)
	case http.MethodPut:
		a.ApiUpdateProduct(w, r, id)
	case http.MethodDelete:
		a.ApiDeleteProduct(w, r, id)
	default:
		writeError(w, http.StatusMethodNotAllowed, "método no permitido")
	}
}

// ApiListProducts devuelve el catálogo completo en formato JSON.
// GET /api/products → 200 []Product
func (a *App) ApiListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := a.products.ListProducts()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error al cargar el catálogo")
		return
	}
	// Si el slice está vacío, devolvemos array vacío en vez de null
	if products == nil {
		products = []models.Product{}
	}
	writeJSON(w, http.StatusOK, products)
}

// ApiGetProduct devuelve un producto específico por su ID.
// GET /api/products/{id} → 200 Product | 404
func (a *App) ApiGetProduct(w http.ResponseWriter, r *http.Request, id string) {
	product, err := a.products.FindProductByID(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "producto no encontrado")
		return
	}
	writeJSON(w, http.StatusOK, product)
}

// ApiCreateProduct lee un JSON del body y crea un nuevo producto.
// POST /api/products → 201 Product | 400 | 403
func (a *App) ApiCreateProduct(w http.ResponseWriter, r *http.Request) {
	if _, ok := a.requireSellerAPI(w, r); !ok {
		return
	}

	// Estructura temporal para leer el JSON del body
	var input struct {
		Name        string  `json:"name"`
		Subtitle    string  `json:"subtitle"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Image       string  `json:"image"`
		Alt         string  `json:"alt"`
		Flavors     string  `json:"flavors"` // Recibido como string separado por comas
		Featured    bool    `json:"featured"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	priceText := strings.TrimSpace(strings.ReplaceAll(
		strings.TrimSpace(fmt.Sprintf("%.2f", input.Price)), ",", "."),
	)

	err := a.products.CreateProduct(
		input.Name,
		input.Subtitle,
		input.Description,
		priceText,
		input.Image,
		input.Alt,
		input.Flavors,
		input.Featured,
	)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Recuperamos el producto recién creado para devolverlo en la respuesta
	products, _ := a.products.ListProducts()
	for i := len(products) - 1; i >= 0; i-- {
		if products[i].Name == input.Name {
			writeJSON(w, http.StatusCreated, products[i])
			return
		}
	}

	writeJSON(w, http.StatusCreated, map[string]string{"ok": "producto creado"})
}

// ApiUpdateProduct lee un JSON del body y actualiza un producto existente.
// PUT /api/products/{id} → 200 Product | 400 | 403 | 404
func (a *App) ApiUpdateProduct(w http.ResponseWriter, r *http.Request, id string) {
	if _, ok := a.requireSellerAPI(w, r); !ok {
		return
	}

	var input struct {
		Name        string  `json:"name"`
		Subtitle    string  `json:"subtitle"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Image       string  `json:"image"`
		Alt         string  `json:"alt"`
		Flavors     string  `json:"flavors"`
		Featured    bool    `json:"featured"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	priceText := strings.TrimSpace(fmt.Sprintf("%.2f", input.Price))

	err := a.products.UpdateProduct(
		id,
		input.Name,
		input.Subtitle,
		input.Description,
		priceText,
		input.Image,
		input.Alt,
		input.Flavors,
		input.Featured,
	)
	if err != nil {
		if strings.Contains(err.Error(), "no encontrado") {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	// Devolvemos el producto actualizado
	product, _ := a.products.FindProductByID(id)
	writeJSON(w, http.StatusOK, product)
}

// ApiDeleteProduct elimina un producto por su ID.
// DELETE /api/products/{id} → 204 | 403 | 404
func (a *App) ApiDeleteProduct(w http.ResponseWriter, r *http.Request, id string) {
	if _, ok := a.requireSellerAPI(w, r); !ok {
		return
	}

	if err := a.products.DeleteProduct(id); err != nil {
		if strings.Contains(err.Error(), "no encontrado") {
			writeError(w, http.StatusNotFound, err.Error())
		} else {
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	// 204 No Content: éxito sin cuerpo de respuesta
	w.WriteHeader(http.StatusNoContent)
}
