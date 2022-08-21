package routes

import (
	"portalnesia.com/api/controllers"
	"portalnesia.com/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetUpUserRouters(r fiber.Router) {
	r.Get("/user", middleware.OnlyLogin, controllers.FindMe)
	r.Get("/user/:id", controllers.FindUser)
}
