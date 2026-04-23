package models

import "time"

// ContactMessage representa la información recibida a través del formulario de contacto.
type ContactMessage struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
