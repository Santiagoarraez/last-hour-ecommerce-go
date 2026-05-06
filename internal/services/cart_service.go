package services

import (
	"errors"

	"lasthour/internal/models"
	"lasthour/internal/storage"
)

// CartService gestiona la persistencia y la lógica del carrito de compras de los usuarios.
type CartService struct {
	carts    *storage.CartStorage
	products *ProductService // Lo necesitamos para obtener detalles de precios y nombres
}

func NewCartService(carts *storage.CartStorage, products *ProductService) *CartService {
	return &CartService{carts: carts, products: products}
}

// AddItem añade un producto al carrito. Si el producto con el mismo sabor ya existe, incrementa cantidad.
func (s *CartService) AddItem(userID, productID string, quantity int, flavors []string, flavorID, flavorName string, price float64, image string) error {
	if quantity < 1 {
		quantity = 1
	}

	carts, err := s.carts.FindAll()
	if err != nil {
		return err
	}

	for cartIndex := range carts {
		if carts[cartIndex].UserID == userID {
			// Buscamos si el mismo sabor ya está en el carrito
			for itemIndex := range carts[cartIndex].Items {
				item := &carts[cartIndex].Items[itemIndex]
				
				// Prioridad al FlavorID para identificar el item
				if flavorID != "" && item.FlavorID == flavorID {
					item.Quantity += quantity
					return s.carts.SaveAll(carts)
				}
				
				// Fallback para compatibilidad con el sistema de "bundles" antiguo
				if flavorID == "" && item.ProductID == productID && s.compareFlavors(item.Flavors, flavors) {
					item.Quantity += quantity
					return s.carts.SaveAll(carts)
				}
			}

			// Si no existe esa combinación, añadimos nuevo item con toda su metadata
			carts[cartIndex].Items = append(carts[cartIndex].Items, models.CartItem{
				ProductID:  productID,
				FlavorID:   flavorID,
				FlavorName: flavorName,
				Price:      price,
				Image:      image,
				Quantity:   quantity,
				Flavors:    flavors,
			})
			return s.carts.SaveAll(carts)
		}
	}

	// Si el usuario no tiene carrito, lo creamos
	carts = append(carts, models.Cart{
		UserID: userID,
		Items: []models.CartItem{{
			ProductID:  productID,
			FlavorID:   flavorID,
			FlavorName: flavorName,
			Price:      price,
			Image:      image,
			Quantity:   quantity,
			Flavors:    flavors,
		}},
	})

	return s.carts.SaveAll(carts)
}

// UpdateQuantity permite modificar la cantidad de un item específico.
func (s *CartService) UpdateQuantity(userID, id string, quantity int) error {
	carts, err := s.carts.FindAll()
	if err != nil {
		return err
	}

	for cartIndex := range carts {
		if carts[cartIndex].UserID == userID {
			for itemIndex := range carts[cartIndex].Items {
				item := &carts[cartIndex].Items[itemIndex]
				// Buscamos coincidencia por FlavorID o por ProductID (fallback)
				if (item.FlavorID != "" && item.FlavorID == id) || (item.FlavorID == "" && item.ProductID == id) {
					if quantity <= 0 {
						return s.RemoveItem(userID, id)
					}
					item.Quantity = quantity
					return s.carts.SaveAll(carts)
				}
			}
		}
	}
	return errors.New("item no encontrado")
}

func (s *CartService) compareFlavors(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// RemoveItem elimina un producto específico del carrito del usuario.
func (s *CartService) RemoveItem(userID, id string) error {
	carts, err := s.carts.FindAll()
	if err != nil {
		return err
	}

	for cartIndex := range carts {
		if carts[cartIndex].UserID == userID {
			var items []models.CartItem
			for _, item := range carts[cartIndex].Items {
				// Comprobamos si coincide el ID enviado con FlavorID o ProductID
				isMatch := (item.FlavorID != "" && item.FlavorID == id) || (item.FlavorID == "" && item.ProductID == id)
				if !isMatch {
					items = append(items, item)
				}
			}
			carts[cartIndex].Items = items
			return s.carts.SaveAll(carts)
		}
	}

	return nil
}

// GetCart obtiene el carrito del usuario convertido a un formato visual útil para plantillas (CartView).
func (s *CartService) GetCart(userID string) (models.CartView, error) {
	carts, err := s.carts.FindAll()
	if err != nil {
		return models.CartView{}, err
	}

	for _, cart := range carts {
		if cart.UserID == userID {
			return s.buildCartView(cart)
		}
	}

	return models.CartView{}, nil
}

// Checkout vacía el carrito del usuario tras realizar un pedido.
func (s *CartService) Checkout(userID string) error {
	carts, err := s.carts.FindAll()
	if err != nil {
		return err
	}

	for cartIndex := range carts {
		if carts[cartIndex].UserID == userID {
			if len(carts[cartIndex].Items) == 0 {
				return errors.New("el carrito esta vacio")
			}
			carts[cartIndex].Items = []models.CartItem{}
			return s.carts.SaveAll(carts)
		}
	}

	return errors.New("el carrito esta vacio")
}

// buildCartView combina datos persistidos con la estructura visual para la web.
func (s *CartService) buildCartView(cart models.Cart) (models.CartView, error) {
	var view models.CartView

	for _, item := range cart.Items {
		var product models.Product

		// Intentamos reconstruir el objeto Product desde la metadata del item (nuevo sistema modular)
		if item.FlavorName != "" {
			product = models.Product{
				ID:    item.FlavorID,
				Name:  item.FlavorName,
				Price: item.Price,
				Image: item.Image,
			}
		} else {
			// Fallback al sistema antiguo buscando en el storage de productos
			p, err := s.products.FindProductByID(item.ProductID)
			if err != nil {
				continue
			}
			product = p
		}

		subtotal := product.Price * float64(item.Quantity)
		view.Items = append(view.Items, models.CartViewItem{
			Product:  product,
			Quantity: item.Quantity,
			Flavors:  item.Flavors,
			Subtotal: subtotal,
		})
		view.Total += subtotal
	}

	return view, nil
}
