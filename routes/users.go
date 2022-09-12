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
	r.Get("/user/:id/followers", user_controllers.FindFollowers)
	r.Get("/user/:id/following", user_controllers.FindFollowings)
	r.Get("/user/:id/followers/pending", user_controllers.FindFollowersPending)
	r.Get("/user/:id/media", user_controllers.FindMedia)
}
