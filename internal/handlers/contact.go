package handlers

import "net/http"

type ContactPageData struct {
	Title   string
	Success bool
	Error   string
}

func (a *App) Contact(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.render(w, "contact.html", ContactPageData{Title: "Contact Us | Vape Store"})
	case http.MethodPost:
		a.processContact(w, r)
	default:
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
	}
}

func (a *App) processContact(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "No se pudo leer el formulario", http.StatusBadRequest)
		return
	}

	err := a.contacts.CreateMessage(
		r.FormValue("name"),
		r.FormValue("email"),
		r.FormValue("message"),
	)
	if err != nil {
		a.render(w, "contact.html", ContactPageData{
			Title: "Contact Us | Vape Store",
			Error: err.Error(),
		})
		return
	}

	a.render(w, "contact.html", ContactPageData{
		Title:   "Contact Us | Vape Store",
		Success: true,
	})
}
