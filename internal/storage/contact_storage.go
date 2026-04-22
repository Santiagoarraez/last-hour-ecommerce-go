package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"lasthour/internal/models"
)

type ContactStorage struct {
	filePath string
}

func NewContactStorage(filePath string) *ContactStorage {
	return &ContactStorage{filePath: filePath}
}

func (s *ContactStorage) Save(message models.ContactMessage) error {
	messages, err := s.FindAll()
	if err != nil {
		return err
	}

	messages = append(messages, message)
	data, err := json.MarshalIndent(messages, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(s.filePath), 0755); err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0644)
}

func (s *ContactStorage) FindAll() ([]models.ContactMessage, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []models.ContactMessage{}, nil
		}
		return nil, err
	}

	if len(data) == 0 {
		return []models.ContactMessage{}, nil
	}

	var messages []models.ContactMessage
	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}
