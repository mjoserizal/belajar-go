package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mjoserizal/belajar-go/controllers/authcontroller"
	"github.com/mjoserizal/belajar-go/controllers/bookcontroller"
	"github.com/mjoserizal/belajar-go/models"
	"github.com/mjoserizal/belajar-go/seed"
)

func main() {
	models.ConnectDatabase()
	seed.Seed(models.DB)
	app := fiber.New()

	api := app.Group("/api")
	book := api.Group("/v1/books")

	book.Get("/", bookcontroller.Index)
	book.Get("/:id", bookcontroller.Show)
	book.Post("/", bookcontroller.Create)
	book.Put("/:id", bookcontroller.Update)
	book.Delete("/:id", bookcontroller.Delete)

	auth := api.Group("/v1")
	auth.Get("/", authcontroller.Index)
	auth.Post("/login", authcontroller.Login)
	auth.Post("/register", authcontroller.Register)
	auth.Post("/logout", authcontroller.Logout)
	auth.Put("/changePassword", authcontroller.ChangePassword)
	auth.Put("/updateUser/:id", authcontroller.UpdateUser)
	auth.Delete("/deleteUser/:id", authcontroller.DeleteUser)
	app.Listen(":8000")
}
