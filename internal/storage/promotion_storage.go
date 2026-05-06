package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"lasthour/internal/models"
)

// PromotionStorage gestiona la persistencia de los packs promocionales.
type PromotionStorage struct {
	filePath string
}

func NewPromotionStorage(filePath string) *PromotionStorage {
	return &PromotionStorage{filePath: filePath}
}

// GetAll lee todas las promociones del archivo JSON.
func (s *PromotionStorage) GetAll() ([]models.Promotion, error) {
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		return []models.Promotion{}, nil
	}

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, err
	}

	var promotions []models.Promotion
	if err := json.Unmarshal(data, &promotions); err != nil {
		return nil, err
	}

	return promotions, nil
}

// GetByID busca una promoción por su ID.
func (s *PromotionStorage) GetByID(id string) (models.Promotion, error) {
	all, err := s.GetAll()
	if err != nil {
		return models.Promotion{}, err
	}

	for _, p := range all {
		if p.ID == id {
			return p, nil
		}
	}

	return models.Promotion{}, errors.New("promocion no encontrada")
}

// Save guarda o actualiza una promoción.
func (s *PromotionStorage) Save(promotion models.Promotion) error {
	all, err := s.GetAll()
	if err != nil {
		return err
	}

	found := false
	for i, p := range all {
		if p.ID == promotion.ID {
			all[i] = promotion
			found = true
			break
		}
	}

	if !found {
		all = append(all, promotion)
	}

	return s.saveAll(all)
}

// Delete elimina una promoción por su ID.
func (s *PromotionStorage) Delete(id string) error {
	all, err := s.GetAll()
	if err != nil {
		return err
	}

	var filtered []models.Promotion
	for _, p := range all {
		if p.ID != id {
			filtered = append(filtered, p)
		}
	}

	if len(filtered) == len(all) {
		return errors.New("promocion no encontrada")
	}

	return s.saveAll(filtered)
}

// saveAll persiste la lista completa de promociones en el archivo JSON.
func (s *PromotionStorage) saveAll(promotions []models.Promotion) error {
	data, err := json.MarshalIndent(promotions, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(s.filePath), 0755); err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0644)
}
