package services

import (
	"errors"

	"lasthour/internal/models"
	"lasthour/internal/storage"
)

type CartService struct {
	carts    *storage.CartStorage
	products *ProductService
}

func NewCartService(carts *storage.CartStorage, products *ProductService) *CartService {
	return &CartService{carts: carts, products: products}
}

func (s *CartService) AddItem(userID, productID string, quantity int) error {
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
			for itemIndex := range carts[cartIndex].Items {
				if carts[cartIndex].Items[itemIndex].ProductID == productID {
					carts[cartIndex].Items[itemIndex].Quantity += quantity
					return s.carts.SaveAll(carts)
				}
			}

			carts[cartIndex].Items = append(carts[cartIndex].Items, models.CartItem{
				ProductID: productID,
				Quantity:  quantity,
			})
			return s.carts.SaveAll(carts)
		}
	}

	carts = append(carts, models.Cart{
		UserID: userID,
		Items:  []models.CartItem{{ProductID: productID, Quantity: quantity}},
	})

	return s.carts.SaveAll(carts)
}

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
			Subtotal: subtotal,
		})
		view.Total += subtotal
	}

	return view, nil
}
