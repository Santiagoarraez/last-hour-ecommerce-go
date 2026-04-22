package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"lasthour/internal/models"
	"lasthour/internal/storage"

	"golang.org/x/crypto/bcrypt"
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		ID:       fmt.Sprintf("user-%d", time.Now().UnixNano()),
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
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
		if user.Email == email {
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err != nil {
				return models.User{}, errors.New("credenciales no validas")
			}
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

func (s *AuthService) UpdateProfile(id, name, email, phone string) (models.User, error) {
	name = strings.TrimSpace(name)
	email = strings.ToLower(strings.TrimSpace(email))
	phone = strings.TrimSpace(phone)

	if name == "" || email == "" {
		return models.User{}, errors.New("nombre y email son obligatorios")
	}

	users, err := s.storage.FindAll()
	if err != nil {
		return models.User{}, err
	}

	for i := range users {
		if users[i].ID == id {
			// Check if email is already taken by another user
			for j := range users {
				if users[j].Email == email && users[j].ID != id {
					return models.User{}, errors.New("ya existe otro usuario con ese correo")
				}
			}

			users[i].Name = name
			users[i].Email = email
			users[i].Phone = phone

			if err := s.storage.SaveAll(users); err != nil {
				return models.User{}, err
			}
			return users[i], nil
		}
	}

	return models.User{}, errors.New("usuario no encontrado")
}
