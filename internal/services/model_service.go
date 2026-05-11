package services

import (
	"errors"
	"strings"

	"lasthour/internal/models"
	"lasthour/internal/storage"
)

// ModelService gestiona la lógica de negocio para los modelos de vape.
type ModelService struct {
	storage       *storage.ModelStorage
	flavorStorage *storage.FlavorStorage
}

func NewModelService(storage *storage.ModelStorage, flavorStorage *storage.FlavorStorage) *ModelService {
	return &ModelService{
		storage:       storage,
		flavorStorage: flavorStorage,
	}
}

// ListModels devuelve la lista completa de modelos.
func (s *ModelService) ListModels() ([]models.VapeModel, error) {
	return s.storage.GetAll()
}

// FindModelByID busca un modelo por su identificador único.
func (s *ModelService) FindModelByID(id string) (models.VapeModel, error) {
	return s.storage.GetByID(id)
}

// CreateModel valida los datos y crea un nuevo modelo de vape.
func (s *ModelService) CreateModel(name, subtitle, description string, price float64) error {
	name = strings.TrimSpace(name)
	subtitle = strings.TrimSpace(subtitle)
	description = strings.TrimSpace(description)

	if name == "" || subtitle == "" || description == "" {
		return errors.New("todos los campos son obligatorios")
	}

	if price <= 0 {
		return errors.New("el precio debe ser mayor a cero")
	}

	id := strings.ToLower(name)
	id = strings.ReplaceAll(id, " ", "-")

	model := models.VapeModel{
		ID:          id,
		Name:        name,
		Subtitle:    subtitle,
		Description: description,
		Price:       price,
	}

	return s.storage.Save(model)
}

// UpdateModel actualiza un modelo y propaga el cambio de nombre a sus sabores.
func (s *ModelService) UpdateModel(id, name, subtitle, description string, price float64) error {
	name = strings.TrimSpace(name)
	subtitle = strings.TrimSpace(subtitle)
	description = strings.TrimSpace(description)

	if name == "" || subtitle == "" || description == "" {
		return errors.New("todos los campos son obligatorios")
	}

	if price <= 0 {
		return errors.New("el precio debe ser mayor a cero")
	}

	model := models.VapeModel{
		ID:          id,
		Name:        name,
		Subtitle:    subtitle,
		Description: description,
		Price:       price,
	}

	// 1. Guardar el modelo actualizado
	if err := s.storage.Save(model); err != nil {
		return err
	}

	// 2. Actualización en cascada: actualizar el nombre del modelo en todos sus sabores
	allFlavors, err := s.flavorStorage.GetAll()
	if err != nil {
		return err
	}

	updatedAny := false
	for i, f := range allFlavors {
		if f.ModelID == id {
			allFlavors[i].ModelName = name
			updatedAny = true
		}
	}

	if updatedAny {
		return s.flavorStorage.SaveAll(allFlavors)
	}

	return nil
}

// DeleteModel elimina un modelo del sistema.
func (s *ModelService) DeleteModel(id string) error {
	return s.storage.Delete(id)
}
