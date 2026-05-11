package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

)


// Cart renderiza la página del carrito del usuario autenticado.
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

	a.render(w, r, "cart.html", map[string]any{
		"Title": "Cart - Last Hour",
		"User":  user,
		"Cart":  cart,
	})
}

// CartAdd añade un producto al carrito. Acepta tanto JSON (AJAX) como form-data (HTML form).
func (a *App) CartAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	user, ok := a.requireUser(w, r)
	if !ok {
		return
	}

	var productID, flavorID, flavorName, image string
	var quantity int
	var price float64
	var flavors []string

	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		var input struct {
			ProductID  string   `json:"product_id"`
			FlavorID   string   `json:"flavor_id"`
			FlavorName string   `json:"flavor_name"`
			Price      float64  `json:"price"`
			Image      string   `json:"image"`
			Quantity   int      `json:"quantity"`
			Flavors    []string `json:"flavors"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "JSON inválido", http.StatusBadRequest)
			return
		}
		productID = input.ProductID
		flavorID = input.FlavorID
		flavorName = input.FlavorName
		price = input.Price
		image = input.Image
		quantity = input.Quantity
		flavors = input.Flavors
	} else {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "No se pudo leer el formulario", http.StatusBadRequest)
			return
		}
		productID = r.FormValue("product_id")
		flavorID = r.FormValue("flavor_id")
		flavorName = r.FormValue("flavor_name")
		if flavorName == "" {
			flavorName = r.FormValue("flavor")
		}
		price, _ = strconv.ParseFloat(r.FormValue("price"), 64)
		image = r.FormValue("image")
		quantity, _ = strconv.Atoi(r.FormValue("quantity"))
		if flavorName != "" {
			flavors = append(flavors, flavorName)
		}
	}

	if quantity <= 0 {
		quantity = 1
	}

	if err := a.carts.AddItem(user.ID, productID, quantity, flavors, flavorID, flavorName, price, image); err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
		return
	}

	if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		cart, _ := a.carts.GetCart(user.ID)
		writeJSON(w, http.StatusOK, map[string]any{"message": "Añadido", "cart": cart})
	} else {
		http.Redirect(w, r, "/cart", http.StatusSeeOther)
	}
}

// CartRemove elimina un item del carrito y redirige de vuelta a /cart.
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

// CartCheckout procesa la orden y redirige al usuario a WhatsApp para finalizar la compra.
// Esta es la funcionalidad estrella de la PEC 2 para gestionar el pedido sin pasarela de pago.
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

	// 1. Construcción del mensaje de WhatsApp con el resumen del pedido.
	// El mensaje incluye nombre del cliente, productos, cantidades y total.
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Hola Last Hour, soy %s. Quisiera realizar el siguiente pedido:\n\n", user.Name))

	for _, item := range cart.Items {
		sb.WriteString(fmt.Sprintf("• %s x%d - %.2f€\n", item.Product.Name, item.Quantity, item.Subtotal))
	}

	sb.WriteString(fmt.Sprintf("\nTotal: %.2f€\n", cart.Total))
	sb.WriteString("\nEspero vuestra confirmación. ¡Gracias!")

	message := sb.String()
	phoneNumber := "34674466462" // Número de contacto de la tienda (ficticio para la PEC)

	// 2. Limpieza del carrito tras "confirmar" el pedido.
	if err := a.carts.Checkout(user.ID); err != nil {
		http.Error(w, "Error al procesar el pedido", http.StatusInternalServerError)
		return
	}

	// 3. Redirección final a la API de WhatsApp.
	// Uso de url.QueryEscape para asegurar que el mensaje se transmita correctamente.
	whatsappURL := fmt.Sprintf("https://wa.me/%s?text=%s", phoneNumber, url.QueryEscape(message))
	http.Redirect(w, r, whatsappURL, http.StatusSeeOther)
}
