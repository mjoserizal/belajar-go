package postcontroller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/mjoserizal/belajar-go/models"
	"gorm.io/gorm"
)

func Index(c *fiber.Ctx) error {
	var posts []models.Post
	models.DB.Find(&posts)

	return c.JSON(posts)
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

	var posts models.Post
	if err := c.BodyParser(&posts); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if err := models.DB.Create(&posts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(posts)
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
