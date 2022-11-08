package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func Database() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}
