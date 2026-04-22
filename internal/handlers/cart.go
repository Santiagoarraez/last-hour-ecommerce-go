package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

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

	cart, err := a.carts.GetCart(user.ID)
	if err != nil || len(cart.Items) == 0 {
		http.Redirect(w, r, "/cart", http.StatusSeeOther)
		return
	}

	// 1. Build the WhatsApp message
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Hola Last Hour, soy %s. Quisiera realizar el siguiente pedido:\n\n", user.Name))

	for _, item := range cart.Items {
		sb.WriteString(fmt.Sprintf("• %s x%d - %.2f€\n", item.Product.Name, item.Quantity, item.Subtotal))
	}

	sb.WriteString(fmt.Sprintf("\nTotal: %.2f€\n", cart.Total))
	sb.WriteString("\nEspero vuestra confirmación. ¡Gracias!")

	message := sb.String()
	phoneNumber := "34674466462" // User's requested number

	// 2. Clear the cart (simulating execution of the order)
	if err := a.carts.Checkout(user.ID); err != nil {
		http.Error(w, "Error al procesar el pedido", http.StatusInternalServerError)
		return
	}

	// 3. Redirect to WhatsApp
	whatsappURL := fmt.Sprintf("https://wa.me/%s?text=%s", phoneNumber, url.QueryEscape(message))
	http.Redirect(w, r, whatsappURL, http.StatusSeeOther)
}
