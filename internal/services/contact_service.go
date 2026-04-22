package services

import (
	"errors"
	"strings"
	"time"

	"lasthour/internal/models"
	"lasthour/internal/storage"
)

type ContactService struct {
	storage *storage.ContactStorage
}

func NewContactService(storage *storage.ContactStorage) *ContactService {
	return &ContactService{storage: storage}
}

func (s *ContactService) CreateMessage(name, email, message string) error {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	message = strings.TrimSpace(message)

	if name == "" || email == "" || message == "" {
		return errors.New("todos los campos son obligatorios")
	}

	contactMessage := models.ContactMessage{
		Name:      name,
		Email:     email,
		Message:   message,
		CreatedAt: time.Now(),
	}

	return s.storage.Save(contactMessage)
}
