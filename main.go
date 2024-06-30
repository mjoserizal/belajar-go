package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/mjoserizal/belajar-go/controllers/authcontroller"
	"github.com/mjoserizal/belajar-go/controllers/bookcontroller"
	"github.com/mjoserizal/belajar-go/models"
	"github.com/mjoserizal/belajar-go/seed"
)

func main() {
	models.ConnectDatabase()
	seed.Seed(models.DB)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))
	api := app.Group("/api")
	book := api.Group("/v1/posts")

	book.Get("/", postcontroller.Index)
	book.Get("/:id", postcontroller.Show)
	book.Post("/", postcontroller.Create)
	book.Put("/:id", postcontroller.Update)
	book.Delete("/:id", postcontroller.Delete)

	auth := api.Group("/v1")
	auth.Get("/", authcontroller.Index)
	auth.Post("/login", authcontroller.Login)
	auth.Post("/register", authcontroller.Register)
	auth.Post("/logout", authcontroller.Logout)
	auth.Put("/changePassword", authcontroller.ChangePassword)
	auth.Put("/updateUser/:id", authcontroller.UpdateUser)
	auth.Delete("/deleteUser/:id", authcontroller.DeleteUser)
	auth.Get("/profile", authcontroller.Profile)
	auth.Put("/updateProfile", authcontroller.UpdateProfile)
	app.Listen(":8000")
}
