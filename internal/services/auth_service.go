package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"lasthour/internal/models"
	"lasthour/internal/storage"
)

type AuthService struct {
	storage *storage.UserStorage
}

func NewAuthService(storage *storage.UserStorage) *AuthService {
	return &AuthService{storage: storage}
}

func (s *AuthService) Register(name, email, password string) (models.User, error) {
	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))
	password = strings.TrimSpace(password)

	if name == "" || email == "" || password == "" {
		return models.User{}, errors.New("todos los campos son obligatorios")
	}

	users, err := s.storage.FindAll()
	if err != nil {
		return models.User{}, err
	}

	for _, user := range users {
		if user.Email == email {
			return models.User{}, errors.New("ya existe un usuario con ese correo")
		}
	}

	user := models.User{
		ID:       fmt.Sprintf("user-%d", time.Now().UnixNano()),
		Name:     name,
		Email:    email,
		Password: password,
		Role:     "customer",
	}

	users = append(users, user)
	if err := s.storage.SaveAll(users); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *AuthService) Login(email, password string) (models.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	password = strings.TrimSpace(password)

	users, err := s.storage.FindAll()
	if err != nil {
		return models.User{}, err
	}

	for _, user := range users {
		if user.Email == email && user.Password == password {
			return user, nil
		}
	}

	return models.User{}, errors.New("credenciales no validas")
}

func (s *AuthService) FindUserByID(id string) (models.User, error) {
	users, err := s.storage.FindAll()
	if err != nil {
		return models.User{}, err
	}

	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}

	return models.User{}, errors.New("usuario no encontrado")
}
