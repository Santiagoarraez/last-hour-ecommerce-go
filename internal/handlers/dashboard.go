package handlers

import "net/http"

// SellerDashboard sirve la página del panel de administración SPA.
// Esta página es una Single Page Application impulsada por la API REST.
func (a *App) SellerDashboard(w http.ResponseWriter, r *http.Request) {
	_, ok := a.requireSeller(w, r)
	if !ok {
		return
	}

	a.render(w, r, "seller_dashboard.html", map[string]any{
		"Title": "Admin Dashboard - Last Hour",
	})
}
