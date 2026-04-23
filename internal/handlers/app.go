package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"

	"lasthour/internal/models"
	"lasthour/internal/services"
)

// App es el núcleo de la aplicación web, el cual orquesta los servicios y el renderizado.
// Contiene las referencias a los servicios de productos, contacto, autenticación y carrito.
type App struct {
	products    *services.ProductService
	contacts    *services.ContactService
	auth        *services.AuthService
	carts       *services.CartService
	templateDir string
}

// NewApp crea una nueva instancia de la aplicación inyectando las dependencias necesarias.
func NewApp(products *services.ProductService, contacts *services.ContactService, auth *services.AuthService, carts *services.CartService, templateDir string) *App {
	return &App{
		products:    products,
		contacts:    contacts,
		auth:        auth,
		carts:       carts,
		templateDir: templateDir,
	}
}

// render se encarga de procesar las plantillas HTML y enviar la respuesta al navegador.
// Siempre incluye el 'layout.html' como base para mantener la cabecera y el pie de página consistentes.
func (a *App) render(w http.ResponseWriter, page string, data any) {
	files := []string{
		filepath.Join(a.templateDir, "layout.html"), // Plantilla base lógica común
		filepath.Join(a.templateDir, page),        // Plantilla específica de la página
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		http.Error(w, "Error cargando la plantilla", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Ejecutamos la plantilla principal 'layout' definida en layout.html
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, "Error generando la respuesta HTML", http.StatusInternalServerError)
	}
}

// currentUser intenta recuperar el usuario actual basado en la cookie de sesión.
// Retorna el usuario y un booleano indicando si se encontró con éxito.
func (a *App) currentUser(r *http.Request) (models.User, bool) {
	cookie, err := r.Cookie("user_id")
	if err != nil || cookie.Value == "" {
		return models.User{}, false
	}

	// Buscamos al usuario según el ID guardado en la cookie
	user, err := a.auth.FindUserByID(cookie.Value)
	if err != nil {
		return models.User{}, false
	}

	return user, true
}

// requireUser es un "middleware" que obliga al usuario a estar autenticado.
// Si no hay sesión, lo redirige al login.
func (a *App) requireUser(w http.ResponseWriter, r *http.Request) (models.User, bool) {
	user, ok := a.currentUser(r)
	if !ok {
		// Redirección si no está logueado
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return models.User{}, false
	}

	return user, true
}

// requireSeller asegura que el usuario autenticado sea además un vendedor.
// Se usa para proteger las rutas de administración de inventario.
func (a *App) requireSeller(w http.ResponseWriter, r *http.Request) (models.User, bool) {
	user, ok := a.requireUser(w, r)
	if !ok {
		return models.User{}, false
	}

	if !user.IsSeller() {
		// Error de prohibido si no tiene el rol 'seller'
		http.Error(w, "Acceso permitido solo para vendedores", http.StatusForbidden)
		return models.User{}, false
	}

	return user, true
}

// setSessionCookie crea una cookie de sesión con el ID del usuario.
// HttpOnly se activa por seguridad para que no sea accesible desde JS.
func setSessionCookie(w http.ResponseWriter, user models.User) {
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    user.ID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// clearSessionCookie borra la cookie de sesión (usado al hacer Logout).
func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Fuerza la expiración inmediata de la cookie
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
