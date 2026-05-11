package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"lasthour/internal/models"
)

// UserStorage gestiona la persistencia de los usuarios en un archivo JSON.
// Garantiza que los perfiles actualizados y los hashes de Bcrypt se guarden correctamente.
type UserStorage struct {
	filePath string
}

func NewUserStorage(filePath string) *UserStorage {
	return &UserStorage{filePath: filePath}
}

// FindAll lee todos los usuarios del archivo JSON.
// Devuelve un slice vacío si el archivo no existe todavía.
func (s *UserStorage) FindAll() ([]models.User, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []models.User{}, nil
		}
		return nil, err
	}

	if len(data) == 0 {
		return []models.User{}, nil
	}

	var users []models.User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// SaveAll persiste la lista completa de usuarios en el archivo JSON.
// Los hashes de bcrypt se guardan tal cual, nunca contraseñas en plano.
func (s *UserStorage) SaveAll(users []models.User) error {
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(s.filePath), 0755); err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0644)
}
