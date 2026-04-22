package handlers

import (
	"net/http"

	"lasthour/internal/models"
)

type AuthPageData struct {
	Title string
	Error string
}

type AccountPageData struct {
	Title string
	User  models.User
}

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

func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	clearSessionCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

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
