package services

import (
	"errors"
	"strings"
	"time"

	"lasthour/internal/models"
	"lasthour/internal/storage"
)

// ContactService gestiona la lógica de mensajes recibidos a través del formulario de contacto.
type ContactService struct {
	storage *storage.ContactStorage
}

func NewContactService(storage *storage.ContactStorage) *ContactService {
	return &ContactService{storage: storage}
}

// CreateMessage valida y guarda un nuevo mensaje de contacto en el sistema.
func (s *ContactService) CreateMessage(name, email, message string) error {
	// Limpiamos los espacios en blanco innecesarios
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	message = strings.TrimSpace(message)

	// Validación de campos obligatorios
	if name == "" || email == "" || message == "" {
		return errors.New("todos los campos son obligatorios")
	}

	contactMessage := models.ContactMessage{
		Name:      name,
		Email:     email,
		Message:   message,
		CreatedAt: time.Now(), // Registramos el momento exacto del envío
	}

	// Persistimos el mensaje a través del almacenamiento
	return s.storage.Save(contactMessage)
}
