package main

import (
	"log"
	"net/http"

	"lasthour/internal/handlers"
	"lasthour/internal/services"
	"lasthour/internal/storage"
)

func main() {
	// 1. Inicialización de la capa de Persistencia (Storage)
	// Definimos las rutas a los archivos JSON que actúan como base de datos
	productStorage := storage.NewProductStorage("data/products.json")
	contactStorage := storage.NewContactStorage("data/messages.json")
	userStorage := storage.NewUserStorage("data/users.json")
	cartStorage := storage.NewCartStorage("data/carts.json")

	// 2. Inicialización de la capa de Negocio (Services) con Inyección de Dependencias
	productService := services.NewProductService(productStorage)
	contactService := services.NewContactService(contactStorage)
	authService := services.NewAuthService(userStorage)
	cartService := services.NewCartService(cartStorage, productService)

	// 3. Inicialización de la capa de Orquestación (Handlers)
	app := handlers.NewApp(productService, contactService, authService, cartService, "templates")
	
	// 4. Configuración del Enrutador (Multiplexor)
	mux := http.NewServeMux()

	// Definición de rutas públicas
	mux.HandleFunc("/", app.Home)
	mux.HandleFunc("/products", app.Products)
	mux.HandleFunc("/products/", app.ProductDetail)
	mux.HandleFunc("/about", app.About)
	mux.HandleFunc("/contact", app.Contact)
	
	// Rutas de autenticación y perfil
	mux.HandleFunc("/login", app.Login)
	mux.HandleFunc("/register", app.Register)
	mux.HandleFunc("/logout", app.Logout)
	mux.HandleFunc("/account", app.Account)
	// PEC 2: Nueva ruta para procesar la actualización del perfil
	mux.HandleFunc("/account/update", app.UpdateAccount)
	
	// Rutas del carrito
	mux.HandleFunc("/cart", app.Cart)
	mux.HandleFunc("/cart/add", app.CartAdd)
	mux.HandleFunc("/cart/remove", app.CartRemove)
	// PEC 2: Nueva ruta para el checkout vía WhatsApp
	mux.HandleFunc("/cart/checkout", app.CartCheckout)
	
	// PEC 3: API REST para Carrito y Cuenta
	mux.HandleFunc("/api/cart", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			app.ApiCart(w, r)
		case http.MethodPost:
			app.ApiCartAdd(w, r)
		case http.MethodPatch:
			app.ApiCartUpdateQuantity(w, r)
		case http.MethodDelete:
			app.ApiCartRemove(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/account", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			app.ApiAccountUpdate(w, r)
		} else {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// PEC 3: Panel de administración SPA (API REST + JavaScript)
	mux.HandleFunc("/seller/dashboard", app.SellerDashboard)

	// PEC 3: Endpoints de la API REST
	// /api/products       → GET (listar) y POST (crear)
	// /api/products/{id}  → GET (obtener), PUT (actualizar) y DELETE (eliminar)
	mux.HandleFunc("/api/products", app.ApiProducts)
	mux.HandleFunc("/api/products/", app.ApiProductByID)

	// 5. Servidor de archivos estáticos (CSS, Imágenes, JavaScript)
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	// PEC 3: Servimos los ficheros JS de la carpeta /js/
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))

	// 6. Lanzamiento del servidor
	log.Println("Servidor web iniciado en http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
