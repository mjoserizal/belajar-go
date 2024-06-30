package postcontroller

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mjoserizal/belajar-go/models"
	"gorm.io/gorm"
)

func Index(c *fiber.Ctx) error {
	userInfo, err := getUserInfoFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	// Check if the user has the required role
	if userInfo.Role != "admin" && userInfo.Role != "user" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Permission denied"})
	}

	// Fetch all posts
	var posts []models.Post
	if err := models.DB.Find(&posts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch posts"})
	}

	return c.JSON(posts)
}

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

func Show(c *fiber.Ctx) error {

	id := c.Params("id")
	var posts models.Post
	if err := models.DB.First(&posts, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"message": "Data tidak ditemukan",
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Data tidak ditemukan",
		})
	}

	return c.JSON(posts)
}

func Create(c *fiber.Ctx) error {
	userInfo, err := getUserInfoFromToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthorized"})
	}

	var post models.Post
	if err := c.BodyParser(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Set the UserID from the authenticated user
	post.UserID = userInfo.ID

	if err := models.DB.Create(&post).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(post)
}

func Update(c *fiber.Ctx) error {

	id := c.Params("id")

	var posts models.Post
	if err := c.BodyParser(&posts); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if models.DB.Where("id = ?", id).Updates(&posts).RowsAffected == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Tidak dapat mengupdate data",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data berhasil diupdate",
	})
}

func Delete(c *fiber.Ctx) error {
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

	var posts models.Post
	if models.DB.Delete(&posts, id).RowsAffected == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"message": "Tidak dapat menghapus data",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data berhasil dihapus",
	})
}
