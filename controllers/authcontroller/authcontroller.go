// authcontroller/authcontroller.go

package authcontroller

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	"github.com/mjoserizal/belajar-go/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"
)

// Index handles fetching all users
func Index(c *fiber.Ctx) error {
	// Retrieve user information from token
	userInfo, err := getUserInfoFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	// Perform authorized actions based on user role or ID
	// For example, fetch all users if the user is admin
	if userInfo.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Permission denied"})
	}

	// If authorized, fetch all users
	var users []models.User
	models.DB.Find(&users)

	return c.JSON(users)
}

// Register handles user registration
func Register(c *fiber.Ctx) error {
	var registerData models.User

	if err := c.BodyParser(&registerData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}

	// Check if username already exists
	var existingUser models.User
	if err := models.DB.Where("username = ?", registerData.Username).First(&existingUser).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "Username already taken"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerData.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to register"})
	}

	user := models.User{
		Username: registerData.Username,
		Name:     registerData.Name,
		Email:    registerData.Email,
		Password: string(hashedPassword),
		Role:     "user", // Set the default role to "user"
	}

	if err := models.DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to register"})
	}

	response := models.ResponseRegister{
		Message: "User Registered",
		Success: true,
		Name:    user.Name,
		Email:   user.Email,
	}

	return c.JSON(response)
}

// Login handles user login and issues a JWT token
func Login(c *fiber.Ctx) error {
	var loginData models.User

	if err := c.BodyParser(&loginData); err != nil {
		fmt.Println("Error parsing login data:", err)
		return c.Status(fiber.StatusBadRequest).JSON(models.ResponseLogin{
			Message: "Invalid request",
			Success: false,
		})
	}

	// Retrieve the user from the database
	var user models.User
	if err := models.DB.Where("username = ?", loginData.Username).First(&user).Error; err != nil {
		fmt.Println("Error retrieving user from database:", err)

		// Handle the case where the user is not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ResponseLogin{
				Message: "Invalid credentials",
				Success: false,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(models.ResponseLogin{
			Message: "User is not registered",
			Success: false,
		})
	}

	// Compare the hashed password with the provided password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		fmt.Println("Error comparing passwords:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(models.ResponseLogin{
			Message: "Invalid credentials",
			Success: false,
		})
	}

	// Generate JWT token
	token, err := generateToken(user.ID, user.Name, user.Role)
	if err != nil {
		fmt.Println("Error generating token:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(models.ResponseLogin{
			Message: "Failed to generate token",
			Success: false,
		})
	}

	return c.JSON(models.ResponseLogin{
		Message: "Login successful",
		Success: true,
		Token:   token,
		Name:    user.Name,
		Email:   user.Email,
		Role:    user.Role, // Include the role in the response
	})
}

// Logout handles user logout
func Logout(c *fiber.Ctx) error {
	// You can implement any logout logic here
	// For example, you may invalidate the JWT token on the client side
	return c.JSON(fiber.Map{"message": "Logout successful"})
}

// ChangePassword handles changing the user's password
func ChangePassword(c *fiber.Ctx) error {
	// Retrieve user information from token
	userInfo, err := getUserInfoFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	var changePasswordData struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}

	if err := c.BodyParser(&changePasswordData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}

	// Retrieve the user from the database
	var user models.User
	if err := models.DB.Where("id = ?", userInfo.ID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to retrieve user"})
	}

	// Compare the hashed old password with the provided old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(changePasswordData.OldPassword)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Invalid credentials"})
	}

	// Hash the new password before saving it to the database
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(changePasswordData.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to change password"})
	}

	// Update the user's password in the database
	if err := models.DB.Model(&user).Update("password", string(hashedNewPassword)).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to change password"})
	}

	return c.JSON(fiber.Map{"message": "Password changed successfully"})
}

// DeleteUser handles the deletion of a user
func DeleteUser(c *fiber.Ctx) error {
	// Extract user ID from the request parameters
	id := c.Params("id")

	// Retrieve user information from token
	userInfo, err := getUserInfoFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	// Check if the user has the 'admin' role
	if userInfo.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Permission denied"})
	}

	// Delete the user with the specified ID
	var user models.User
	if models.DB.Delete(&user, id).RowsAffected == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": "Unable to delete user"})
	}

	return c.JSON(fiber.Map{"message": "User deleted successfully"})
}

// Helper function to generate JWT token
func generateToken(userID uint, name, role string) (string, error) {
	// Define the expiration time for the token (e.g., 1 hour)
	expirationTime := time.Now().Add(time.Hour * 1)

	// Create the standard claims
	standardClaims := jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Issuer:    "your-issuer", // Replace with your issuer
		Subject:   fmt.Sprint(userID),
	}

	// Create the custom claims
	customClaims := UserInfo{
		ID:   userID,
		Name: name,
		Role: role,
	}

	// Combine standard claims and custom claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"standard": standardClaims,
		"custom":   customClaims,
	})

	// Sign the token with a secret key and get the complete encoded token as a string
	secretKey := []byte("ey51616516") // Replace with your secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// getUserInfoFromToken retrieves user information from the JWT token
func getUserInfoFromToken(c *fiber.Ctx) (*UserInfo, error) {
	// Get the Authorization header from the request
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("Authorization header missing")
	}

	// Extract the token from the Authorization header
	tokenString := authHeader[len("Bearer "):]

	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Provide the secret key for token validation
		return []byte("ey51616516"), nil
	})
	if err != nil {
		return nil, err
	}

	// Extract claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("Invalid token claims")
	}

	// Extract user information from the custom claims
	userInfo := &UserInfo{
		ID:   uint(claims["custom"].(map[string]interface{})["ID"].(float64)),
		Name: claims["custom"].(map[string]interface{})["Name"].(string),
		Role: claims["custom"].(map[string]interface{})["Role"].(string),
	}

	return userInfo, nil
}

// UserInfo struct to hold user information
type UserInfo struct {
	ID   uint
	Name string
	Role string
}

func UpdateUser(c *fiber.Ctx) error {
	// Extract user ID from the request parameters
	id := c.Params("id")

	// Retrieve user information from token
	userInfo, err := getUserInfoFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	// Check if the user has the 'admin' role
	if userInfo.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Permission denied"})
	}

	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}

	// Allow updates to the Username, Name, and Role fields
	updateData := map[string]interface{}{
		"Username": user.Username,
		"Name":     user.Name,
		"Role":     user.Role, // Allow updating the role
	}

	// Update the user's data in the database
	if models.DB.Model(&models.User{}).Where("id = ?", id).Updates(updateData).RowsAffected == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Permission Denied"})
	}

	return c.JSON(fiber.Map{"message": "User data updated successfully"})
}

// Profile handles fetching user profile
func Profile(c *fiber.Ctx) error {
	// Retrieve user information from token
	userInfo, err := getUserInfoFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	// Retrieve the user from the database
	var user models.User
	if err := models.DB.Where("id = ?", userInfo.ID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to retrieve user"})
	}

	// Exclude sensitive information like password before returning
	user.Password = ""

	return c.JSON(user)
}

// UpdateProfile handles updating user profile
func UpdateProfile(c *fiber.Ctx) error {
	// Retrieve user information from token
	userInfo, err := getUserInfoFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	// Retrieve the user from the database
	var user models.User
	if err := models.DB.Where("id = ?", userInfo.ID).First(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to retrieve user"})
	}

	// Parse and decode the request body into user object
	var updatedUser models.User
	if err := c.BodyParser(&updatedUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request"})
	}

	// Update the user's profile information
	user.Username = updatedUser.Username
	user.Name = updatedUser.Name

	// Save the updated user profile to the database
	if err := models.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to update profile"})
	}

	// Return the updated user profile
	return c.JSON(user)
}
