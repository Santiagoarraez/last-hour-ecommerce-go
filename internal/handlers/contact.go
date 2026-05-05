package handlers

import "net/http"


// Contact gestiona tanto la visualización del formulario de contacto (GET)
// como la recepción del mensaje (POST).
func (a *App) Contact(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Mostramos el formulario vacío
		a.render(w, r, "contact.html", map[string]any{"Title": "Contact Us | Vape Store"})
	case http.MethodPost:
		// Procesamos los datos enviados
		a.processContact(w, r)
	default:
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
	}
}

// processContact extrae los datos del formulario y llama al servicio de mensajes.
func (a *App) processContact(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "No se pudo leer el formulario", http.StatusBadRequest)
		return
	}

	// Enviamos la información al servicio de contacto para su almacenamiento
	err := a.contacts.CreateMessage(
		r.FormValue("name"),
		r.FormValue("email"),
		r.FormValue("message"),
	)
	if err != nil {
		// Si hay error en la lógica de negocio, lo notificamos en la misma página
		a.render(w, r, "contact.html", map[string]any{
			"Title": "Contact Us | Vape Store",
			"Error": err.Error(),
		})
		return
	}

	// Si todo sale bien, mostramos un mensaje de confirmación
	a.render(w, r, "contact.html", map[string]any{
		"Title":   "Contact Us | Vape Store",
		"Success": true,
	})
}
