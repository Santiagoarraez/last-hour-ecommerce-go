package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"

	"lasthour/internal/models"
	"lasthour/internal/services"
)

type App struct {
	products    *services.ProductService
	contacts    *services.ContactService
	auth        *services.AuthService
	carts       *services.CartService
	templateDir string
}

func NewApp(products *services.ProductService, contacts *services.ContactService, auth *services.AuthService, carts *services.CartService, templateDir string) *App {
	return &App{
		products:    products,
		contacts:    contacts,
		auth:        auth,
		carts:       carts,
		templateDir: templateDir,
	}
}

func (a *App) render(w http.ResponseWriter, page string, data any) {
	files := []string{
		filepath.Join(a.templateDir, "layout.html"),
		filepath.Join(a.templateDir, page),
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "Error cargando la plantilla", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, "Error generando la respuesta HTML", http.StatusInternalServerError)
	}
}

func (a *App) currentUser(r *http.Request) (models.User, bool) {
	cookie, err := r.Cookie("user_id")
	if err != nil || cookie.Value == "" {
		return models.User{}, false
	}

	user, err := a.auth.FindUserByID(cookie.Value)
	if err != nil {
		return models.User{}, false
	}

	return user, true
}

func (a *App) requireUser(w http.ResponseWriter, r *http.Request) (models.User, bool) {
	user, ok := a.currentUser(r)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return models.User{}, false
	}

	return user, true
}

func (a *App) requireSeller(w http.ResponseWriter, r *http.Request) (models.User, bool) {
	user, ok := a.requireUser(w, r)
	if !ok {
		return models.User{}, false
	}

	if !user.IsSeller() {
		http.Error(w, "Acceso permitido solo para vendedores", http.StatusForbidden)
		return models.User{}, false
	}

	return user, true
}

func setSessionCookie(w http.ResponseWriter, user models.User) {
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    user.ID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
