package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"portalnesia.com/api/database"
	"portalnesia.com/api/models"
	"portalnesia.com/api/response"
)

func FindUsers(c *fiber.Ctx) {
	db := database.DB
	var users []models.User

	db.Limit(15).Find(&users)

	// fmt.Printf("%+v\n", users)

	c.Status(fiber.StatusOK).JSON(fiber.Map{"data": users})
}

func FindUser(c *fiber.Ctx) error {
	db := database.DB
	id := c.Params("id", "")

	fmt.Print(id)
	var user models.User

	if err := db.First(&user, "user_login = ?", id).Error; err != nil {
		return response.NotFound("user", id, "username")
	}

	return response.Response(user).Send(c)
}
