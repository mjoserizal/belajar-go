// models/user.go

package models

type User struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role" gorm:"default:user"`
}
type Response struct {
	Name    string `json:"name"`
	Role    string `json:"role"`
	Message string `json:"message"`
	Success bool   `json:"success"`
	Token   string `json:"token,omitempty"`
	// Add other fields as needed
}
