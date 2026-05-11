package main

import (
	"log"
	"net/http"
	"strings"

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
	modelStorage := storage.NewModelStorage("data/models.json")
	flavorStorage := storage.NewFlavorStorage("data/flavors.json")
	promotionStorage := storage.NewPromotionStorage("data/promotions.json")

	// 2. Inicialización de la capa de Negocio (Services) con Inyección de Dependencias
	productService := services.NewProductService(productStorage)
	contactService := services.NewContactService(contactStorage)
	authService := services.NewAuthService(userStorage)
	cartService := services.NewCartService(cartStorage, productService)
	modelService := services.NewModelService(modelStorage, flavorStorage)
	flavorService := services.NewFlavorService(flavorStorage)
	promotionService := services.NewPromotionService(promotionStorage)
	sessionService := services.NewSessionService()

	// PEC 3: Migración automática de nombres de modelos en sabores al arrancar
	migrateFlavorModelNames(modelService, flavorService, flavorStorage)

	// 3. Inicialización de la capa de Orquestación (Handlers)
	// NewApp también carga y cachea todas las plantillas HTML al arrancar
	app, err := handlers.NewApp(productService, contactService, authService, cartService, modelService, flavorService, promotionService, sessionService, "templates")
	if err != nil {
		log.Fatalf("Error inicializando la aplicación: %v", err)
	}

	// 4. Configuración del Enrutador (Multiplexor)
	mux := http.NewServeMux()

	// Definición de rutas públicas
	mux.HandleFunc("/", app.Home)
	mux.HandleFunc("/products", app.Products)
	mux.HandleFunc("/products/flavor/", app.FlavorDetail)
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

	// Endpoints para Modelos
	mux.HandleFunc("/api/models", app.ApiModels)
	mux.HandleFunc("/api/models/", func(w http.ResponseWriter, r *http.Request) {
		// Diferenciar entre /api/models/{id} y /api/models/{id}/flavors
		if strings.HasSuffix(r.URL.Path, "/flavors") {
			app.ApiListFlavorsByModel(w, r)
		} else {
			app.ApiModelByID(w, r)
		}
	})

	// Endpoints para Sabores
	mux.HandleFunc("/api/flavors", app.ApiFlavors)
	mux.HandleFunc("/api/flavors/", app.ApiFlavorByID)

	// Endpoints para Promociones
	mux.HandleFunc("/api/promotions", app.ApiPromotions)
	mux.HandleFunc("/api/promotions/", app.ApiPromotionByID)

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

// migrateFlavorModelNames sincroniza los nombres de modelos en los sabores si están vacíos.
func migrateFlavorModelNames(ms *services.ModelService, fs *services.FlavorService, fst *storage.FlavorStorage) {
	flavors, err := fs.ListFlavors()
	if err != nil {
		log.Printf("Error migración sabores: %v", err)
		return
	}

	models, err := ms.ListModels()
	if err != nil {
		log.Printf("Error migración modelos: %v", err)
		return
	}

	// Mapa de ID -> Nombre para búsqueda rápida
	modelMap := make(map[string]string)
	for _, m := range models {
		modelMap[m.ID] = m.Name
	}

	updated := false
	for i := range flavors {
		// Si el nombre del modelo está vacío, intentamos recuperarlo del mapa
		if flavors[i].ModelName == "" {
			if name, ok := modelMap[flavors[i].ModelID]; ok {
				flavors[i].ModelName = name
				updated = true
			}
		}
	}

	if updated {
		if err := fst.SaveAll(flavors); err != nil {
			log.Printf("Error guardando sabores migrados: %v", err)
		} else {
			log.Println("Migración exitosa: ModelName actualizado en sabores.")
		}
	}
}
