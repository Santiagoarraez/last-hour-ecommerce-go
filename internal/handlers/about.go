package handlers

import "net/http"


// About maneja la visualización de la página de información de la empresa.
func (a *App) About(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	a.render(w, r, "about.html", map[string]any{"Title": "About Us | Vape Store"})
}
