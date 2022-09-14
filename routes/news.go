package routes

import (
	"github.com/gofiber/fiber/v2"
	news_controllers "portalnesia.com/api/controllers/news"
	"portalnesia.com/api/middleware"
)

func SetUpNewsRouters(r fiber.Router) {
	r.Use(middleware.OnlyInternal)
	r.Get("/news", news_controllers.FindNews)
	r.Get("/news/:source/:title", news_controllers.FindOneNews)
}
