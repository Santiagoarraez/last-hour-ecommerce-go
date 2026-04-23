package models

// User representa la estructura de datos para los usuarios del sistema.
// Se ha añadido el campo "Phone" para permitir el perfil de usuario completo.
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// IsSeller verifica si el usuario tiene el rol de vendedor.
func (u User) IsSeller() bool {
	return u.Role == "seller"
}
