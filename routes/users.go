package routes

import (
	"portalnesia.com/api/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetUpUserRouters(r fiber.Router) {
	// r.GET("/user", controllers.FindUsers)
	r.Get("/user/:id", controllers.FindUser)
}
