package handlers

import (
	"net/http"

	"lasthour/internal/models"
)

// AuthPageData contiene los datos necesarios para renderizar las páginas de Login y Registro.
type AuthPageData struct {
	Title string
	Error string
}

// AccountPageData contiene los datos para la página de gestión de cuenta (Perfil).
// Estructura creada para soportar la visualización y edición del perfil.
type AccountPageData struct {
	Title   string
	User    models.User
	Success string
	Error   string
}

// Login maneja la visualización del formulario y el proceso de inicio de sesión.
func (a *App) Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.render(w, "login.html", AuthPageData{Title: "Login - Last Hour"})
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, "No se pudo leer el formulario", http.StatusBadRequest)
			return
		}

		user, err := a.auth.Login(r.FormValue("email"), r.FormValue("password"))
		if err != nil {
			a.render(w, "login.html", AuthPageData{Title: "Login - Last Hour", Error: err.Error()})
			return
		}

		setSessionCookie(w, user)
		http.Redirect(w, r, "/account", http.StatusSeeOther)
	default:
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
	}
}

// Register maneja el registro de nuevos usuarios.
func (a *App) Register(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.render(w, "register.html", AuthPageData{Title: "Register - Last Hour"})
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, "No se pudo leer el formulario", http.StatusBadRequest)
			return
		}

		user, err := a.auth.Register(r.FormValue("name"), r.FormValue("email"), r.FormValue("password"))
		if err != nil {
			a.render(w, "register.html", AuthPageData{Title: "Register - Last Hour", Error: err.Error()})
			return
		}

		setSessionCookie(w, user)
		http.Redirect(w, r, "/account", http.StatusSeeOther)
	default:
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
	}
}

// Logout cierra la sesión del usuario.
func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	clearSessionCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Account renderiza la página de perfil del usuario.
// Nueva página para que el usuario gestione sus datos.
func (a *App) Account(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	user, ok := a.requireUser(w, r)
	if !ok {
		return
	}

	a.render(w, "account.html", AccountPageData{
		Title: "Account - Last Hour",
		User:  user,
	})
}

// UpdateAccount procesa la actualización de los datos del perfil.
// Handler añadido para procesar el formulario de edición de perfil.
func (a *App) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	user, ok := a.requireUser(w, r)
	if !ok {
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "No se pudo leer el formulario", http.StatusBadRequest)
		return
	}

	// Llamada al servicio para actualizar los datos (incluyendo el nuevo campo Phone)
	updatedUser, err := a.auth.UpdateProfile(
		user.ID,
		r.FormValue("name"),
		r.FormValue("email"),
		r.FormValue("phone"),
	)

	if err != nil {
		a.render(w, "account.html", AccountPageData{
			Title: "Account - Last Hour",
			User:  user,
			Error: err.Error(),
		})
		return
	}

	// Actualizamos la cookie de sesión con la nueva información del usuario
	setSessionCookie(w, updatedUser)

	a.render(w, "account.html", AccountPageData{
		Title:   "Account - Last Hour",
		User:    updatedUser,
		Success: "Perfil actualizado correctamente",
	})
}
