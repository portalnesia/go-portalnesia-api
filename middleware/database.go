package middleware

import (
	"github.com/gofiber/fiber/v2"
	"portalnesia.com/api/config"
	"portalnesia.com/api/models"
)

func Database() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		config.ChangeDatabase(false)
		c.Locals("ctx", &models.CtxDefaultValue)
		return c.Next()
	}
}
