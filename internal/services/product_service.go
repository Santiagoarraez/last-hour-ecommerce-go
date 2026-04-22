package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"lasthour/internal/models"
	"lasthour/internal/storage"
)

type ProductService struct {
	storage *storage.ProductStorage
}

func NewProductService(storage *storage.ProductStorage) *ProductService {
	return &ProductService{storage: storage}
}

func (s *ProductService) ListProducts() ([]models.Product, error) {
	return s.storage.FindAll()
}

func (s *ProductService) ListFeaturedProducts() ([]models.Product, error) {
	products, err := s.storage.FindAll()
	if err != nil {
		return nil, err
	}

	var featured []models.Product
	for _, product := range products {
		if product.Featured {
			featured = append(featured, product)
		}
	}

	return featured, nil
}

func (s *ProductService) FindProductByID(id string) (models.Product, error) {
	products, err := s.storage.FindAll()
	if err != nil {
		return models.Product{}, err
	}

	for _, product := range products {
		if product.ID == id {
			return product, nil
		}
	}

	return models.Product{}, errors.New("producto no encontrado")
}

func (s *ProductService) CreateProduct(name, subtitle, description, priceText, image, alt, flavorsText string, featured bool) error {
	price, err := strconv.ParseFloat(strings.TrimSpace(priceText), 64)
	if err != nil || price <= 0 {
		return errors.New("el precio debe ser un numero positivo")
	}

	product := models.Product{
		ID:          buildProductID(name),
		Name:        strings.TrimSpace(name),
		Subtitle:    strings.TrimSpace(subtitle),
		Description: strings.TrimSpace(description),
		Price:       price,
		Image:       strings.TrimSpace(image),
		Alt:         strings.TrimSpace(alt),
		Flavors:     splitFlavors(flavorsText),
		Featured:    featured,
	}

	if product.Name == "" || product.Subtitle == "" || product.Description == "" || product.Image == "" {
		return errors.New("nombre, subtitulo, descripcion e imagen son obligatorios")
	}

	products, err := s.storage.FindAll()
	if err != nil {
		return err
	}

	for _, existing := range products {
		if existing.ID == product.ID {
			product.ID = fmt.Sprintf("%s-%d", product.ID, time.Now().Unix())
			break
		}
	}

	products = append(products, product)
	return s.storage.SaveAll(products)
}

func (s *ProductService) UpdateProduct(id, name, subtitle, description, priceText, image, alt, flavorsText string, featured bool) error {
	price, err := strconv.ParseFloat(strings.TrimSpace(priceText), 64)
	if err != nil || price <= 0 {
		return errors.New("el precio debe ser un numero positivo")
	}

	products, err := s.storage.FindAll()
	if err != nil {
		return err
	}

	for index := range products {
		if products[index].ID == id {
			products[index].Name = strings.TrimSpace(name)
			products[index].Subtitle = strings.TrimSpace(subtitle)
			products[index].Description = strings.TrimSpace(description)
			products[index].Price = price
			products[index].Image = strings.TrimSpace(image)
			products[index].Alt = strings.TrimSpace(alt)
			products[index].Flavors = splitFlavors(flavorsText)
			products[index].Featured = featured
			return s.storage.SaveAll(products)
		}
	}

	return errors.New("producto no encontrado")
}

func (s *ProductService) DeleteProduct(id string) error {
	products, err := s.storage.FindAll()
	if err != nil {
		return err
	}

	var filtered []models.Product
	for _, product := range products {
		if product.ID != id {
			filtered = append(filtered, product)
		}
	}

	if len(filtered) == len(products) {
		return errors.New("producto no encontrado")
	}

	return s.storage.SaveAll(filtered)
}

func splitFlavors(value string) []string {
	parts := strings.Split(value, ",")
	var flavors []string
	for _, part := range parts {
		flavor := strings.TrimSpace(part)
		if flavor != "" {
			flavors = append(flavors, flavor)
		}
	}
	return flavors
}

func buildProductID(name string) string {
	value := strings.ToLower(strings.TrimSpace(name))
	value = strings.ReplaceAll(value, " ", "-")
	value = strings.ReplaceAll(value, "_", "-")

	var builder strings.Builder
	lastDash := false
	for _, char := range value {
		if char >= 'a' && char <= 'z' || char >= '0' && char <= '9' {
			builder.WriteRune(char)
			lastDash = false
			continue
		}

		if char == '-' && !lastDash {
			builder.WriteRune(char)
			lastDash = true
		}
	}

	result := strings.Trim(builder.String(), "-")
	if result == "" {
		return fmt.Sprintf("product-%d", time.Now().Unix())
	}
	return result
}
