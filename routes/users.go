package routes

import (
	user_controllers "portalnesia.com/api/controllers/users"
	"portalnesia.com/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetUpUserRouters(r fiber.Router) {
	r.Get("/user", middleware.OnlyLogin, user_controllers.FindMe)
	r.Get("/user/list", middleware.OnlyInternal, user_controllers.ListUsername)
	r.Get("/user/:id/followers/pending", middleware.OnlySpecificScope([]string{"user"}), middleware.OnlyLogin, user_controllers.FindFollowersPending)
	r.Get("/user/:id/followers", middleware.OnlySpecificScope([]string{"user"}), user_controllers.FindFollowers)
	r.Get("/user/:id/following", middleware.OnlySpecificScope([]string{"user"}), user_controllers.FindFollowings)
	r.Get("/user/:id/media", middleware.OnlySpecificScope([]string{"user", "files"}), user_controllers.FindMedia)
	r.Get("/user/:id", middleware.OnlySpecificScope([]string{"user"}), user_controllers.FindUser)

	r.Post("/user/:id", middleware.OnlySpecificScope([]string{"user-write"}), middleware.OnlyLogin, user_controllers.UploadPhotoProfile)
}
