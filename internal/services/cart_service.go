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

// AddItem añade un producto al carrito o incrementa su cantidad si ya existe.
func (s *CartService) AddItem(userID, productID string, quantity int) error {
	if quantity < 1 {
		quantity = 1
	}

	// Verificamos que el producto exista realmente antes de añadirlo
	if _, err := s.products.FindProductByID(productID); err != nil {
		return err
	}

	carts, err := s.carts.FindAll()
	if err != nil {
		return err
	}

	// Buscamos si el usuario ya tiene un carrito iniciado
	for cartIndex := range carts {
		if carts[cartIndex].UserID == userID {
			// Buscamos si el producto ya está en el carrito
			for itemIndex := range carts[cartIndex].Items {
				if carts[cartIndex].Items[itemIndex].ProductID == productID {
					// Solo incrementamos la cantidad
					carts[cartIndex].Items[itemIndex].Quantity += quantity
					return s.carts.SaveAll(carts)
				}
			}

			// Si el producto no estaba, lo añadimos a la lista de items
			carts[cartIndex].Items = append(carts[cartIndex].Items, models.CartItem{
				ProductID: productID,
				Quantity:  quantity,
			})
			return s.carts.SaveAll(carts)
		}
	}

	// Si el usuario no tenía carrito, creamos uno nuevo
	carts = append(carts, models.Cart{
		UserID: userID,
		Items:  []models.CartItem{{ProductID: productID, Quantity: quantity}},
	})

	return s.carts.SaveAll(carts)
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
			// Filtramos los items manteniendo todos menos el solicitado
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
			// Limpiamos los items del carrito
			carts[cartIndex].Items = []models.CartItem{}
			return s.carts.SaveAll(carts)
		}
	}

	return errors.New("el carrito esta vacio")
}

// buildCartView es una función auxiliar que combina los IDs del carrito con los datos reales de los productos.
// Calcula subtotales y el total acumulado.
func (s *CartService) buildCartView(cart models.Cart) (models.CartView, error) {
	var view models.CartView

	for _, item := range cart.Items {
		// Obtenemos los detalles del producto (nombre, precio, etc.)
		product, err := s.products.FindProductByID(item.ProductID)
		if err != nil {
			continue // Si un producto se borró del catálogo, lo ignoramos en el carrito
		}

		subtotal := product.Price * float64(item.Quantity)
		view.Items = append(view.Items, models.CartViewItem{
			Product:  product,
			Quantity: item.Quantity,
			Subtotal: subtotal,
		})
		view.Total += subtotal
	}

	return view, nil
}
