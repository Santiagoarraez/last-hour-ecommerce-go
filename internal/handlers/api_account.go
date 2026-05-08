package handlers

import (
	"encoding/json"
	"net/http"
)

// ApiAccountUpdate actualiza el perfil del usuario autenticado.
// PUT /api/account
func (a *App) ApiAccountUpdate(w http.ResponseWriter, r *http.Request) {
	user, ok := a.currentUser(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "debes iniciar sesión")
		return
	}

	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	updatedUser, err := a.auth.UpdateProfile(user.ID, input.Name, input.Email, input.Phone)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// No es necesario actualizar la sesión porque el token no cambia
	// y sigue apuntando al mismo ID de usuario.

	writeJSON(w, http.StatusOK, map[string]any{
		"message": "Perfil actualizado",
		"user":    updatedUser,
	})
}
