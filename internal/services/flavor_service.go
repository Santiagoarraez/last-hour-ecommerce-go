package services

import (
	"errors"
	"strings"

	"lasthour/internal/models"
	"lasthour/internal/storage"
)

// FlavorService gestiona la lógica de negocio para los sabores de los vapes.
type FlavorService struct {
	storage *storage.FlavorStorage
}

func NewFlavorService(storage *storage.FlavorStorage) *FlavorService {
	return &FlavorService{storage: storage}
}

// ListFlavors devuelve todos los sabores registrados.
func (s *FlavorService) ListFlavors() ([]models.Flavor, error) {
	return s.storage.GetAll()
}

// ListFlavorsByModel devuelve los sabores asociados a un modelo específico.
func (s *FlavorService) ListFlavorsByModel(modelID string) ([]models.Flavor, error) {
	return s.storage.GetByModelID(modelID)
}

// FindFlavorByID busca un sabor por su ID único.
func (s *FlavorService) FindFlavorByID(id string) (models.Flavor, error) {
	return s.storage.GetByID(id)
}

// CreateFlavor valida y registra un nuevo sabor para un modelo.
func (s *FlavorService) CreateFlavor(modelID, modelName, name, image string) error {
	modelID = strings.TrimSpace(modelID)
	name = strings.TrimSpace(name)

	if modelID == "" || name == "" {
		return errors.New("el ID del modelo y el nombre del sabor son obligatorios")
	}

	// Generar ID: modelid-nombre-del-sabor
	id := strings.ToLower(modelID + "-" + name)
	id = strings.ReplaceAll(id, " ", "-")

	flavor := models.Flavor{
		ID:        id,
		ModelID:   modelID,
		ModelName: modelName,
		Name:      name,
		Image:     image,
	}

	return s.storage.Save(flavor)
}

// UpdateFlavor actualiza los datos de un sabor existente.
func (s *FlavorService) UpdateFlavor(id, name, image string) error {
	flavor, err := s.storage.GetByID(id)
	if err != nil {
		return err
	}

	name = strings.TrimSpace(name)
	if name != "" {
		flavor.Name = name
	}
	
	// Solo actualizamos la imagen si se proporciona una nueva (base64)
	if image != "" {
		flavor.Image = image
	}

	return s.storage.Save(flavor)
}

// SetStock marca un sabor como agotado o disponible.
func (s *FlavorService) SetStock(id string, outOfStock bool) error {
	flavor, err := s.storage.GetByID(id)
	if err != nil {
		return err
	}
	flavor.OutOfStock = outOfStock
	return s.storage.Save(flavor)
}

// DeleteFlavor elimina un sabor del sistema.
func (s *FlavorService) DeleteFlavor(id string) error {
	return s.storage.Delete(id)
}
