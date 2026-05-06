package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"lasthour/internal/models"
)

// ModelStorage se encarga de gestionar la persistencia de los modelos de vapes.
type ModelStorage struct {
	filePath string
}

func NewModelStorage(filePath string) *ModelStorage {
	return &ModelStorage{filePath: filePath}
}

// GetAll lee todos los modelos del archivo JSON.
func (s *ModelStorage) GetAll() ([]models.VapeModel, error) {
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		return []models.VapeModel{}, nil
	}

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, err
	}

	var vapeModels []models.VapeModel
	if err := json.Unmarshal(data, &vapeModels); err != nil {
		return nil, err
	}

	return vapeModels, nil
}

// GetByID busca un modelo por su ID.
func (s *ModelStorage) GetByID(id string) (models.VapeModel, error) {
	vapeModels, err := s.GetAll()
	if err != nil {
		return models.VapeModel{}, err
	}

	for _, m := range vapeModels {
		if m.ID == id {
			return m, nil
		}
	}

	return models.VapeModel{}, errors.New("modelo no encontrado")
}

// Save guarda un modelo (creación o actualización).
func (s *ModelStorage) Save(model models.VapeModel) error {
	vapeModels, err := s.GetAll()
	if err != nil {
		return err
	}

	found := false
	for i, m := range vapeModels {
		if m.ID == model.ID {
			vapeModels[i] = model
			found = true
			break
		}
	}

	if !found {
		vapeModels = append(vapeModels, model)
	}

	return s.saveAll(vapeModels)
}

// Delete elimina un modelo por su ID.
func (s *ModelStorage) Delete(id string) error {
	vapeModels, err := s.GetAll()
	if err != nil {
		return err
	}

	var filtered []models.VapeModel
	for _, m := range vapeModels {
		if m.ID != id {
			filtered = append(filtered, m)
		}
	}

	if len(filtered) == len(vapeModels) {
		return errors.New("modelo no encontrado")
	}

	return s.saveAll(filtered)
}

// saveAll persiste la lista completa de modelos en el archivo JSON.
func (s *ModelStorage) saveAll(vapeModels []models.VapeModel) error {
	data, err := json.MarshalIndent(vapeModels, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(s.filePath), 0755); err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0644)
}
