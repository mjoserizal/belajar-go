// seed/seed.go

package seed

import (
	"log"

	"github.com/mjoserizal/belajar-go/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Seed initializes the database with seed data
func Seed(db *gorm.DB) {
	// Check if the user table is already populated
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		log.Println("User table already seeded.")
		return
	}

	// Create an admin user
	hashedAdminPassword, err := bcrypt.GenerateFromPassword([]byte("adminpassword"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	adminUser := models.User{
		Username: "adminuser",
		Password: string(hashedAdminPassword),
		Name:     "Admin",
		Role:     "admin",
	}

	// Insert the admin user into the user table
	if err := db.Create(&adminUser).Error; err != nil {
		log.Fatal(err)
	}

	log.Println("Admin user seeded successfully.")
}
