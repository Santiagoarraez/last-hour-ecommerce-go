package handlers

import (
	"net/http"
	"strings"

	"lasthour/internal/models"
)

type SellerProductsPageData struct {
	Title    string
	User     models.User
	Products []models.Product
	Error    string
}

type SellerProductFormData struct {
	Title   string
	User    models.User
	Product models.Product
	Action  string
	Error   string
}

func (a *App) SellerProducts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/seller/products" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	user, ok := a.requireSeller(w, r)
	if !ok {
		return
	}

	products, err := a.products.ListProducts()
	if err != nil {
		http.Error(w, "No se pudo cargar el catalogo", http.StatusInternalServerError)
		return
	}

	a.render(w, "seller_products.html", SellerProductsPageData{
		Title:    "Seller Products - Last Hour",
		User:     user,
		Products: products,
	})
}

func (a *App) SellerProductNew(w http.ResponseWriter, r *http.Request) {
	user, ok := a.requireSeller(w, r)
	if !ok {
		return
	}

	switch r.Method {
	case http.MethodGet:
		a.render(w, "seller_product_form.html", SellerProductFormData{
			Title:  "New Product - Last Hour",
			User:   user,
			Action: "/seller/products/new",
			Product: models.Product{
				Image: "/assets/images/hqd-catalog-new.png",
			},
		})
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, "No se pudo leer el formulario", http.StatusBadRequest)
			return
		}

		err := a.products.CreateProduct(
			r.FormValue("name"),
			r.FormValue("subtitle"),
			r.FormValue("description"),
			r.FormValue("price"),
			r.FormValue("image"),
			r.FormValue("alt"),
			r.FormValue("flavors"),
			r.FormValue("featured") == "on",
		)
		if err != nil {
			a.render(w, "seller_product_form.html", SellerProductFormData{
				Title:  "New Product - Last Hour",
				User:   user,
				Action: "/seller/products/new",
				Error:  err.Error(),
			})
			return
		}

		http.Redirect(w, r, "/seller/products", http.StatusSeeOther)
	default:
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
	}
}

func (a *App) SellerProductEdit(w http.ResponseWriter, r *http.Request) {
	user, ok := a.requireSeller(w, r)
	if !ok {
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/seller/products/edit/")
	if id == "" || strings.Contains(id, "/") {
		http.NotFound(w, r)
		return
	}

	product, err := a.products.FindProductByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		a.render(w, "seller_product_form.html", SellerProductFormData{
			Title:   "Edit Product - Last Hour",
			User:    user,
			Product: product,
			Action:  "/seller/products/edit/" + id,
		})
	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, "No se pudo leer el formulario", http.StatusBadRequest)
			return
		}

		err := a.products.UpdateProduct(
			id,
			r.FormValue("name"),
			r.FormValue("subtitle"),
			r.FormValue("description"),
			r.FormValue("price"),
			r.FormValue("image"),
			r.FormValue("alt"),
			r.FormValue("flavors"),
			r.FormValue("featured") == "on",
		)
		if err != nil {
			a.render(w, "seller_product_form.html", SellerProductFormData{
				Title:   "Edit Product - Last Hour",
				User:    user,
				Product: product,
				Action:  "/seller/products/edit/" + id,
				Error:   err.Error(),
			})
			return
		}

		http.Redirect(w, r, "/seller/products", http.StatusSeeOther)
	default:
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
	}
}

func (a *App) SellerProductDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}

	if _, ok := a.requireSeller(w, r); !ok {
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/seller/products/delete/")
	if id == "" || strings.Contains(id, "/") {
		http.NotFound(w, r)
		return
	}

	if err := a.products.DeleteProduct(id); err != nil {
		http.Error(w, "No se pudo eliminar el producto", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/seller/products", http.StatusSeeOther)
}
