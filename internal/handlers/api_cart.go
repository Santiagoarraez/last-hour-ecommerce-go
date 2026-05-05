package handlers

import (
	"encoding/json"
	"net/http"
)

// ApiCart devuelve el estado actual del carrito del usuario autenticado.
// GET /api/cart
func (a *App) ApiCart(w http.ResponseWriter, r *http.Request) {
	user, ok := a.currentUser(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "debes iniciar sesión")
		return
	}

	cart, err := a.carts.GetCart(user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "error al obtener el carrito")
		return
	}

	writeJSON(w, http.StatusOK, cart)
}

// ApiCartAdd añade un producto al carrito vía JSON.
// POST /api/cart
func (a *App) ApiCartAdd(w http.ResponseWriter, r *http.Request) {
	user, ok := a.currentUser(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "debes iniciar sesión")
		return
	}

	var input struct {
		ProductID string   `json:"product_id"`
		Quantity  int      `json:"quantity"`
		Flavors   []string `json:"flavors"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if input.Quantity <= 0 {
		input.Quantity = 1
	}

	if err := a.carts.AddItem(user.ID, input.ProductID, input.Quantity, input.Flavors); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	cart, _ := a.carts.GetCart(user.ID)
	writeJSON(w, http.StatusOK, map[string]any{
		"message": "Producto añadido",
		"cart":    cart,
	})
}

// ApiCartUpdateQuantity modifica la cantidad de un item.
// PATCH /api/cart
func (a *App) ApiCartUpdateQuantity(w http.ResponseWriter, r *http.Request) {
	user, ok := a.currentUser(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "debes iniciar sesión")
		return
	}

	var input struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if err := a.carts.UpdateQuantity(user.ID, input.ProductID, input.Quantity); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	cart, _ := a.carts.GetCart(user.ID)
	writeJSON(w, http.StatusOK, map[string]any{
		"message": "Cantidad actualizada",
		"cart":    cart,
	})
}

// ApiCartRemove elimina un producto del carrito.
// DELETE /api/cart
func (a *App) ApiCartRemove(w http.ResponseWriter, r *http.Request) {
	user, ok := a.currentUser(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "debes iniciar sesión")
		return
	}

	productID := r.URL.Query().Get("product_id")
	if productID == "" {
		writeError(w, http.StatusBadRequest, "product_id es requerido")
		return
	}

	if err := a.carts.RemoveItem(user.ID, productID); err != nil {
		writeError(w, http.StatusInternalServerError, "error al eliminar producto")
		return
	}

	cart, _ := a.carts.GetCart(user.ID)
	writeJSON(w, http.StatusOK, map[string]any{
		"message": "Producto eliminado",
		"cart":    cart,
	})
}
