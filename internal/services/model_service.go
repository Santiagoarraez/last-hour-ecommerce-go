package services

import (
	"errors"
	"strconv"
	"strings"

	"lasthour/internal/models"
	"lasthour/internal/storage"
)

// ModelService gestiona la lógica de negocio para los modelos de vape.
type ModelService struct {
	storage *storage.ModelStorage
}

func NewModelService(storage *storage.ModelStorage) *ModelService {
	return &ModelService{storage: storage}
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
func (s *ModelService) CreateModel(name, subtitle, description, priceText string) error {
	name = strings.TrimSpace(name)
	subtitle = strings.TrimSpace(subtitle)
	description = strings.TrimSpace(description)
	priceText = strings.TrimSpace(priceText)

	if name == "" || subtitle == "" || description == "" || priceText == "" {
		return errors.New("todos los campos son obligatorios")
	}

	price, err := strconv.ParseFloat(priceText, 64)
	if err != nil || price <= 0 {
		return errors.New("precio invalido")
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

// DeleteModel elimina un modelo del sistema.
func (s *ModelService) DeleteModel(id string) error {
	return s.storage.Delete(id)
}
