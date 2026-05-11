package handlers

import (
	"fmt"
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
	vapeModels  *services.ModelService
	flavors     *services.FlavorService
	promotions  *services.PromotionService
	sessions    *services.SessionService
	templateDir string
	// templates almacena las plantillas pre-compiladas al arrancar la app.
	// La clave es el nombre del archivo de página (ej: "products.html").
	templates   map[string]*template.Template
}

// NewApp crea una nueva instancia de la aplicación inyectando las dependencias necesarias.
// También carga y cachea todas las plantillas HTML para mejorar el rendimiento.
func NewApp(
	products *services.ProductService,
	contacts *services.ContactService,
	auth *services.AuthService,
	carts *services.CartService,
	vapeModels *services.ModelService,
	flavors *services.FlavorService,
	promotions *services.PromotionService,
	sessions *services.SessionService,
	templateDir string,
) (*App, error) {
	app := &App{
		products:    products,
		contacts:    contacts,
		auth:        auth,
		carts:       carts,
		vapeModels:  vapeModels,
		flavors:     flavors,
		promotions:  promotions,
		sessions:    sessions,
		templateDir: templateDir,
	}

	if err := app.LoadTemplates(); err != nil {
		return nil, fmt.Errorf("error cargando plantillas: %w", err)
	}

	return app, nil
}

// templateFuncs registra funciones disponibles en todos los templates.
// safeURL permite que las data URLs (base64) pasen sin ser sanitizadas por html/template.
var templateFuncs = template.FuncMap{
	"safeURL": func(s string) template.URL { return template.URL(s) },
}

// LoadTemplates recorre el directorio de plantillas, parsea cada archivo HTML
// junto con layout.html y almacena el resultado en el mapa de caché.
// Se excluye el propio layout.html ya que es una dependencia, no una página independiente.
func (a *App) LoadTemplates() error {
	layoutPath := filepath.Join(a.templateDir, "layout.html")

	pages, err := filepath.Glob(filepath.Join(a.templateDir, "*.html"))
	if err != nil {
		return err
	}

	a.templates = make(map[string]*template.Template, len(pages))

	for _, page := range pages {
		name := filepath.Base(page)

		if name == "layout.html" {
			continue
		}

		// Usamos New+Funcs antes de ParseFiles para que las funciones estén disponibles
		tmpl, err := template.New(filepath.Base(layoutPath)).Funcs(templateFuncs).ParseFiles(layoutPath, page)
		if err != nil {
			return fmt.Errorf("error parseando plantilla '%s': %w", name, err)
		}

		a.templates[name] = tmpl
	}

	return nil
}

// render busca la plantilla cacheada por nombre y ejecuta el template "layout".
// Si la plantilla no se encuentra en el mapa devuelve un error 500.
func (a *App) render(w http.ResponseWriter, r *http.Request, page string, data any) {
	tmpl, ok := a.templates[page]
	if !ok {
		http.Error(w, fmt.Sprintf("plantilla '%s' no encontrada", page), http.StatusInternalServerError)
		return
	}

	// Recuperamos el usuario para inyectarlo siempre en todas las vistas
	user, _ := a.currentUser(r)

	finalData := make(map[string]any)
	finalData["User"] = user

	// Si data es un mapa, volcamos sus claves al mapa final
	if m, ok := data.(map[string]any); ok {
		for k, v := range m {
			finalData[k] = v
		}
	} else if data != nil {
		finalData["Data"] = data
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "layout", finalData); err != nil {
		http.Error(w, "Error generando la respuesta HTML", http.StatusInternalServerError)
	}
}

// currentUser intenta recuperar el usuario actual basado en la cookie de sesión.
// La cookie almacena un token opaco; SessionService lo resuelve al userID real.
func (a *App) currentUser(r *http.Request) (models.User, bool) {
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		return models.User{}, false
	}

	// Resolvemos el token al userID (verifica existencia y expiración)
	userID, err := a.sessions.GetUserID(cookie.Value)
	if err != nil {
		return models.User{}, false
	}

	user, err := a.auth.FindUserByID(userID)
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

// setSessionCookie almacena el token opaco de sesión en una cookie HttpOnly.
// El token no contiene información del usuario, solo es una clave de lookup.
func setSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400, // 24 horas en segundos
	})
}

// clearSessionCookie borra la cookie de sesión (usado al hacer Logout).
func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Fuerza la expiración inmediata de la cookie
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
