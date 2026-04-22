package handlers

import "net/http"

type AboutPageData struct {
	Title string
}

func (a *App) About(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/about" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	a.render(w, "about.html", AboutPageData{Title: "About Us | Vape Store"})
}
