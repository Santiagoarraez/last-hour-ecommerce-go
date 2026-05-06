package services

import (
	"errors"
	"strconv"
	"strings"

	"lasthour/internal/models"
	"lasthour/internal/storage"
)

// PromotionService gestiona la lógica de negocio para los packs promocionales.
type PromotionService struct {
	storage *storage.PromotionStorage
}

func NewPromotionService(storage *storage.PromotionStorage) *PromotionService {
	return &PromotionService{storage: storage}
}

// ListPromotions devuelve todas las promociones disponibles.
func (s *PromotionService) ListPromotions() ([]models.Promotion, error) {
	return s.storage.GetAll()
}

// FindPromotionByID busca una promoción por su ID.
func (s *PromotionService) FindPromotionByID(id string) (models.Promotion, error) {
	return s.storage.GetByID(id)
}

// CreatePromotion valida y registra un nuevo pack promocional.
func (s *PromotionService) CreatePromotion(name, description, priceText, unitsText, image string, items []models.PromotionItem) error {
	name = strings.TrimSpace(name)
	priceText = strings.TrimSpace(priceText)
	unitsText = strings.TrimSpace(unitsText)

	if name == "" || priceText == "" || unitsText == "" {
		return errors.New("nombre, precio y unidades son campos obligatorios")
	}

	price, err := strconv.ParseFloat(priceText, 64)
	if err != nil || price <= 0 {
		return errors.New("precio invalido")
	}

	units, err := strconv.Atoi(unitsText)
	if err != nil || units <= 0 {
		return errors.New("unidades invalidas")
	}

	id := strings.ToLower(name)
	id = strings.ReplaceAll(id, " ", "-")

	promotion := models.Promotion{
		ID:          id,
		Name:        name,
		Description: description,
		Price:       price,
		Image:       image,
		Units:       units,
		Items:       items,
	}

	return s.storage.Save(promotion)
}

// UpdatePromotion actualiza los datos de una promoción existente.
func (s *PromotionService) UpdatePromotion(promotion models.Promotion) error {
	return s.storage.Save(promotion)
}

// DeletePromotion elimina una promoción del catálogo.
func (s *PromotionService) DeletePromotion(id string) error {
	return s.storage.Delete(id)
}
