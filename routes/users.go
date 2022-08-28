package routes

import (
	user_controllers "portalnesia.com/api/controllers/users"
	"portalnesia.com/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetUpUserRouters(r fiber.Router) {
	r.Get("/user", middleware.OnlyLogin, user_controllers.FindMe)
	r.Get("/user/list", middleware.OnlyInternal, user_controllers.ListUsername)
	r.Get("/user/:id", user_controllers.FindUser)

}
