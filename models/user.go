// models/user.go

package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `gorm:"unique"`
	Password string `json:"password"`
	Role     string `json:"role" gorm:"default:user"`
	Posts    []Post
}
type ResponseLogin struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	Message string `json:"message"`
	Success bool   `json:"success"`
	Token   string `json:"token,omitempty"`
	// Add other fields as needed
}
type ResponseRegister struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}
