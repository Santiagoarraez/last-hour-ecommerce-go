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

// AddItem añade un producto al carrito. Si el producto con los mismos sabores ya existe, incrementa cantidad.
func (s *CartService) AddItem(userID, productID string, quantity int, flavors []string) error {
	if quantity < 1 {
		quantity = 1
	}

	if _, err := s.products.FindProductByID(productID); err != nil {
		return err
	}

	carts, err := s.carts.FindAll()
	if err != nil {
		return err
	}

	for cartIndex := range carts {
		if carts[cartIndex].UserID == userID {
			// Buscamos si el mismo producto con mismos sabores ya está
			for itemIndex := range carts[cartIndex].Items {
				item := &carts[cartIndex].Items[itemIndex]
				if item.ProductID == productID && s.compareFlavors(item.Flavors, flavors) {
					item.Quantity += quantity
					return s.carts.SaveAll(carts)
				}
			}

			// Si no existe esa combinación, añadimos nuevo item
			carts[cartIndex].Items = append(carts[cartIndex].Items, models.CartItem{
				ProductID: productID,
				Quantity:  quantity,
				Flavors:   flavors,
			})
			return s.carts.SaveAll(carts)
		}
	}

	carts = append(carts, models.Cart{
		UserID: userID,
		Items:  []models.CartItem{{ProductID: productID, Quantity: quantity, Flavors: flavors}},
	})

	return s.carts.SaveAll(carts)
}

// UpdateQuantity permite modificar la cantidad de un item específico.
func (s *CartService) UpdateQuantity(userID, productID string, quantity int) error {
	carts, err := s.carts.FindAll()
	if err != nil {
		return err
	}

	for cartIndex := range carts {
		if carts[cartIndex].UserID == userID {
			for itemIndex := range carts[cartIndex].Items {
				if carts[cartIndex].Items[itemIndex].ProductID == productID {
					if quantity <= 0 {
						// Si la cantidad es 0 o menos, eliminamos el item
						return s.RemoveItem(userID, productID)
					}
					carts[cartIndex].Items[itemIndex].Quantity = quantity
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
func (s *CartService) RemoveItem(userID, productID string) error {
	carts, err := s.carts.FindAll()
	if err != nil {
		return err
	}

	for cartIndex := range carts {
		if carts[cartIndex].UserID == userID {
			var items []models.CartItem
			for _, item := range carts[cartIndex].Items {
				if item.ProductID != productID {
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

// buildCartView combina IDs con detalles reales.
func (s *CartService) buildCartView(cart models.Cart) (models.CartView, error) {
	var view models.CartView

	for _, item := range cart.Items {
		product, err := s.products.FindProductByID(item.ProductID)
		if err != nil {
			continue
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
