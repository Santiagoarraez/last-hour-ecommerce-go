package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"lasthour/internal/models"
)

// FlavorStorage gestiona la persistencia de los sabores asociados a los modelos de vape.
type FlavorStorage struct {
	filePath string
}

func NewFlavorStorage(filePath string) *FlavorStorage {
	return &FlavorStorage{filePath: filePath}
}

// GetAll lee todos los sabores del archivo JSON.
func (s *FlavorStorage) GetAll() ([]models.Flavor, error) {
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		return []models.Flavor{}, nil
	}

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, err
	}

	var flavors []models.Flavor
	if err := json.Unmarshal(data, &flavors); err != nil {
		return nil, err
	}

	return flavors, nil
}

// GetByModelID devuelve los sabores filtrados por el ID del modelo.
func (s *FlavorStorage) GetByModelID(modelID string) ([]models.Flavor, error) {
	all, err := s.GetAll()
	if err != nil {
		return nil, err
	}

	var filtered []models.Flavor
	for _, f := range all {
		if f.ModelID == modelID {
			filtered = append(filtered, f)
		}
	}

	return filtered, nil
}

// GetByID busca un sabor por su ID único.
func (s *FlavorStorage) GetByID(id string) (models.Flavor, error) {
	all, err := s.GetAll()
	if err != nil {
		return models.Flavor{}, err
	}

	for _, f := range all {
		if f.ID == id {
			return f, nil
		}
	}

	return models.Flavor{}, errors.New("sabor no encontrado")
}

// Save guarda o actualiza un sabor.
func (s *FlavorStorage) Save(flavor models.Flavor) error {
	all, err := s.GetAll()
	if err != nil {
		return err
	}

	found := false
	for i, f := range all {
		if f.ID == flavor.ID {
			all[i] = flavor
			found = true
			break
		}
	}

	if !found {
		all = append(all, flavor)
	}

	return s.saveAll(all)
}

// Delete elimina un sabor por su ID.
func (s *FlavorStorage) Delete(id string) error {
	all, err := s.GetAll()
	if err != nil {
		return err
	}

	var filtered []models.Flavor
	for _, f := range all {
		if f.ID != id {
			filtered = append(filtered, f)
		}
	}

	if len(filtered) == len(all) {
		return errors.New("sabor no encontrado")
	}

	return s.saveAll(filtered)
}

// saveAll persiste la lista completa de sabores en el archivo JSON.
func (s *FlavorStorage) saveAll(flavors []models.Flavor) error {
	data, err := json.MarshalIndent(flavors, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(s.filePath), 0755); err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0644)
}
