package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"lasthour/internal/models"
)

// CartStorage almacena los carritos de compra activos de los usuarios en un archivo JSON.
type CartStorage struct {
	filePath string
}

func NewCartStorage(filePath string) *CartStorage {
	return &CartStorage{filePath: filePath}
}

func (s *CartStorage) FindAll() ([]models.Cart, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []models.Cart{}, nil
		}
		return nil, err
	}

	if len(data) == 0 {
		return []models.Cart{}, nil
	}

	var carts []models.Cart
	if err := json.Unmarshal(data, &carts); err != nil {
		return nil, err
	}

	return carts, nil
}

func (s *CartStorage) SaveAll(carts []models.Cart) error {
	data, err := json.MarshalIndent(carts, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(s.filePath), 0755); err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0644)
}
