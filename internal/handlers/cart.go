package handlers

import (
	"net/http"
	"strconv"

	"lasthour/internal/models"
)

type CartPageData struct {
	Title   string
	User    models.User
	Cart    models.CartView
	Success string
	Error   string
}

func (a *App) Cart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	user, ok := a.requireUser(w, r)
	if !ok {
		return
	}

	cart, err := a.carts.GetCart(user.ID)
	if err != nil {
		http.Error(w, "No se pudo cargar el carrito", http.StatusInternalServerError)
		return
	}

	a.render(w, "cart.html", CartPageData{
		Title: "Cart - Last Hour",
		User:  user,
		Cart:  cart,
	})
}

func (a *App) CartAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	user, ok := a.requireUser(w, r)
	if !ok {
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "No se pudo leer el formulario", http.StatusBadRequest)
		return
	}

	quantity, _ := strconv.Atoi(r.FormValue("quantity"))
	if err := a.carts.AddItem(user.ID, r.FormValue("product_id"), quantity); err != nil {
		http.Error(w, "No se pudo agregar el producto", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

func (a *App) CartRemove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	user, ok := a.requireUser(w, r)
	if !ok {
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "No se pudo leer el formulario", http.StatusBadRequest)
		return
	}

	if err := a.carts.RemoveItem(user.ID, r.FormValue("product_id")); err != nil {
		http.Error(w, "No se pudo eliminar el producto", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

func (a *App) CartCheckout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	user, ok := a.requireUser(w, r)
	if !ok {
		return
	}

	if err := a.carts.Checkout(user.ID); err != nil {
		cart, _ := a.carts.GetCart(user.ID)
		a.render(w, "cart.html", CartPageData{
			Title: "Cart - Last Hour",
			User:  user,
			Cart:  cart,
			Error: err.Error(),
		})
		return
	}

	a.render(w, "cart.html", CartPageData{
		Title:   "Cart - Last Hour",
		User:    user,
		Success: "Order processed by the server. Your cart is now empty.",
	})
}
