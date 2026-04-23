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

// ProductService contiene la lógica de negocio para gestionar el catálogo de vapes.
type ProductService struct {
	storage *storage.ProductStorage
}

func NewProductService(storage *storage.ProductStorage) *ProductService {
	return &ProductService{storage: storage}
}

// ListProducts devuelve todos los productos disponibles en el sistema.
func (s *ProductService) ListProducts() ([]models.Product, error) {
	return s.storage.FindAll()
}

// ListFeaturedProducts filtra y devuelve solo los productos marcados como destacados (para la home).
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

// FindProductByID busca un producto específico por su identificador único.
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

// CreateProduct valida y añade un nuevo producto al catálogo.
func (s *ProductService) CreateProduct(name, subtitle, description, priceText, image, alt, flavorsText string, featured bool) error {
	// Conversión de precio de texto a decimal (float64)
	price, err := strconv.ParseFloat(strings.TrimSpace(priceText), 64)
	if err != nil || price <= 0 {
		return errors.New("el precio debe ser un numero positivo")
	}

	product := models.Product{
		ID:          buildProductID(name), // Generamos un ID amigable (slug) basado en el nombre
		Name:        strings.TrimSpace(name),
		Subtitle:    strings.TrimSpace(subtitle),
		Description: strings.TrimSpace(description),
		Price:       price,
		Image:       strings.TrimSpace(image),
		Alt:         strings.TrimSpace(alt),
		Flavors:     splitFlavors(flavorsText), // Convertimos la lista de sabores de texto a un array
		Featured:    featured,
	}

	// Validación básica de campos requeridos
	if product.Name == "" || product.Subtitle == "" || product.Description == "" || product.Image == "" {
		return errors.New("nombre, subtitulo, descripcion e imagen son obligatorios")
	}

	products, err := s.storage.FindAll()
	if err != nil {
		return err
	}

	// Evitamos duplicidad de IDs
	for _, existing := range products {
		if existing.ID == product.ID {
			product.ID = fmt.Sprintf("%s-%d", product.ID, time.Now().Unix())
			break
		}
	}

	products = append(products, product)
	return s.storage.SaveAll(products)
}

// UpdateProduct modifica un producto existente identificado por ID.
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
			// Actualización de todos los datos del producto
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

// DeleteProduct elimina un producto del catálogo permanentemente.
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

// splitFlavors es una función auxiliar que convierte una cadena separada por comas en un slice de strings.
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

// buildProductID genera una cadena amigable para URLs (slug) a partir del nombre del producto.
func buildProductID(name string) string {
	value := strings.ToLower(strings.TrimSpace(name))
	value = strings.ReplaceAll(value, " ", "-")
	value = strings.ReplaceAll(value, "_", "-")

	var builder strings.Builder
	lastDash := false
	for _, char := range value {
		// Solo permitimos caracteres alfanuméricos y guiones
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
		// Fallback por seguridad con timestamp
		return fmt.Sprintf("product-%d", time.Now().Unix())
	}
	return result
}
