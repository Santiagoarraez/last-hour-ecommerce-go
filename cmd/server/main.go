package main

import (
	"log"
	"net/http"

	"lasthour/internal/handlers"
	"lasthour/internal/services"
	"lasthour/internal/storage"
)

func main() {
	productStorage := storage.NewProductStorage("data/products.json")
	contactStorage := storage.NewContactStorage("data/messages.json")
	userStorage := storage.NewUserStorage("data/users.json")
	cartStorage := storage.NewCartStorage("data/carts.json")

	productService := services.NewProductService(productStorage)
	contactService := services.NewContactService(contactStorage)
	authService := services.NewAuthService(userStorage)
	cartService := services.NewCartService(cartStorage, productService)

	app := handlers.NewApp(productService, contactService, authService, cartService, "templates")
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.Home)
	mux.HandleFunc("/products", app.Products)
	mux.HandleFunc("/products/", app.ProductDetail)
	mux.HandleFunc("/about", app.About)
	mux.HandleFunc("/contact", app.Contact)
	mux.HandleFunc("/login", app.Login)
	mux.HandleFunc("/register", app.Register)
	mux.HandleFunc("/logout", app.Logout)
	mux.HandleFunc("/account", app.Account)
	mux.HandleFunc("/account/update", app.UpdateAccount)
	mux.HandleFunc("/cart", app.Cart)
	mux.HandleFunc("/cart/add", app.CartAdd)
	mux.HandleFunc("/cart/remove", app.CartRemove)
	mux.HandleFunc("/cart/checkout", app.CartCheckout)
	mux.HandleFunc("/seller/products", app.SellerProducts)
	mux.HandleFunc("/seller/products/new", app.SellerProductNew)
	mux.HandleFunc("/seller/products/edit/", app.SellerProductEdit)
	mux.HandleFunc("/seller/products/delete/", app.SellerProductDelete)

	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	log.Println("Servidor web iniciado en http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
