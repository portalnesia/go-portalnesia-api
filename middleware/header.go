package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mssola/user_agent"
	"portalnesia.com/api/config"
)

func setUpIP(c *fiber.Ctx) {
	ip := ""
	if config.NODE_ENV == "production" {
		ip = c.Get("x-local-api", "")
		if ip == "" {
			ip = c.Get("cf-connecting-ip", "")
		}
		if ip == "" {
			ip = c.IP()
		}
	} else {
		ip = c.IP()
	}
	c.Locals("ip", ip)
}

func setUpBrowser(c *fiber.Ctx) {
	browser := user_agent.New(c.Get("user-agent"))
	br, ver := browser.Browser()
	browserString := fmt.Sprintf("%s, %s %s", browser.OS(), br, ver)

	c.Locals("browser", browser)
	c.Locals("browserStr", browserString)
}

func Header() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Set("Vary", "Accept-Encoding")
		setUpIP(c)
		setUpBrowser(c)
		return c.Next()
	}
}
