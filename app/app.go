package app

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"portalnesia.com/api/config"
	"portalnesia.com/api/models"
	"portalnesia.com/api/routes"
)

func Initialization() {
	config.SetupConfig()
	config.SetupFirebase()
	models.SetupDB()
	if os.Getenv("NODE_ENV") == "production" {
		models.SetupDebugDB()
	}
	config.ChangeDatabase(false)
}

func NewApp() *fiber.App {
	// Setup gofiber
	r := routes.SetupRouters()
	return r
}
