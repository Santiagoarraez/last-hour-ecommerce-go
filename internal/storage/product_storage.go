package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"lasthour/internal/models"
)

// ProductStorage se encarga de leer y escribir los productos en el archivo JSON del sistema.
type ProductStorage struct {
	filePath string
}

func NewProductStorage(filePath string) *ProductStorage {
	return &ProductStorage{filePath: filePath}
}

// FindAll lee y deserializa todos los productos del archivo JSON.
func (s *ProductStorage) FindAll() ([]models.Product, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, err
	}

	var products []models.Product
	if err := json.Unmarshal(data, &products); err != nil {
		return nil, err
	}

	return products, nil
}

// SaveAll serializa y persiste la lista completa de productos en el archivo JSON.
func (s *ProductStorage) SaveAll(products []models.Product) error {
	data, err := json.MarshalIndent(products, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(s.filePath), 0755); err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0644)
}
