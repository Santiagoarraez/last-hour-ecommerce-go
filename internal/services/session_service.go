package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

// session representa una entrada en el mapa de sesiones activas.
type session struct {
	userID    string
	expiresAt time.Time
}

// SessionService gestiona los tokens de sesión en memoria.
// Utiliza un mapa protegido con RWMutex para ser seguro ante accesos concurrentes.
type SessionService struct {
	mu       sync.RWMutex
	sessions map[string]session
}

// NewSessionService crea una nueva instancia de SessionService e inicia
// una goroutine de limpieza que elimina sesiones expiradas cada hora.
func NewSessionService() *SessionService {
	svc := &SessionService{
		sessions: make(map[string]session),
	}

	// Limpieza periódica en background para evitar fugas de memoria
	go svc.cleanupLoop()

	return svc
}

// CreateSession genera un token aleatorio seguro de 32 bytes (64 caracteres hex),
// lo asocia al userID con una expiración de 24 horas y lo devuelve.
func (s *SessionService) CreateSession(userID string) (string, error) {
	// Generamos 32 bytes de aleatoriedad criptográfica
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}

	token := hex.EncodeToString(raw)

	s.mu.Lock()
	s.sessions[token] = session{
		userID:    userID,
		expiresAt: time.Now().Add(24 * time.Hour),
	}
	s.mu.Unlock()

	return token, nil
}

// GetUserID resuelve un token a un userID.
// Devuelve error si el token no existe o ha expirado.
// Las sesiones expiradas se eliminan de forma lazy en este punto.
func (s *SessionService) GetUserID(token string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sess, ok := s.sessions[token]
	if !ok {
		return "", errors.New("sesión no encontrada")
	}

	if time.Now().After(sess.expiresAt) {
		// Eliminamos la sesión expirada de forma lazy
		delete(s.sessions, token)
		return "", errors.New("sesión expirada")
	}

	return sess.userID, nil
}

// DeleteSession elimina un token del mapa, invalidando la sesión.
// Se llama al hacer logout.
func (s *SessionService) DeleteSession(token string) {
	s.mu.Lock()
	delete(s.sessions, token)
	s.mu.Unlock()
}

// cleanupLoop es una goroutine que borra sesiones expiradas cada hora
// para evitar que el mapa crezca indefinidamente.
func (s *SessionService) cleanupLoop() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		s.mu.Lock()
		for token, sess := range s.sessions {
			if now.After(sess.expiresAt) {
				delete(s.sessions, token)
			}
		}
		s.mu.Unlock()
	}
}
