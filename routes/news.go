package routes

import (
	"github.com/gofiber/fiber/v2"
	"portalnesia.com/api/controllers"
)

func SetUpNewsRouters(r fiber.Router) {
	r.Get("/news", controllers.FindNews)
	r.Get("/news/:source/:title", controllers.FindOneNews)
}
