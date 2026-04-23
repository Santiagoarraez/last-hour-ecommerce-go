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
	
	// Rutas de administración (Vendedor)
	mux.HandleFunc("/seller/products", app.SellerProducts)
	mux.HandleFunc("/seller/products/new", app.SellerProductNew)
	mux.HandleFunc("/seller/products/edit/", app.SellerProductEdit)
	mux.HandleFunc("/seller/products/delete/", app.SellerProductDelete)

	// 5. Servidor de archivos estáticos (CSS, Imágenes)
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// 6. Lanzamiento del servidor
	log.Println("Servidor web iniciado en http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
