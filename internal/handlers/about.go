package handlers

import "net/http"

// AboutPageData define los datos para la página estática "Sobre Nosotros".
type AboutPageData struct {
	Title string
}

// About maneja la visualización de la página de información de la empresa.
func (a *App) About(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	a.render(w, "about.html", AboutPageData{Title: "About Us | Vape Store"})
}
