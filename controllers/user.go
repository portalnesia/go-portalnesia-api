package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"portalnesia.com/api/config"
	"portalnesia.com/api/models"
	"portalnesia.com/api/response"
	util "portalnesia.com/api/utils"
)

func FindMe(c *fiber.Ctx) error {
	db := config.DB
	ctx := c.Locals("ctx").(*models.Context)

	withEmail := ctx != nil && (ctx.IsWeb || ctx.Client != nil && ctx.Client.Scope != nil && util.CheckScope(*ctx.Client.Scope, []string{"email"}))
	fmt.Println(ctx.IsInternal)
	users := ctx.ToUserModels(db, withEmail)
	return response.Response(users).Send(c)
}

func FindUser(c *fiber.Ctx) error {
	db := config.DB
	id := c.Params("id", "")

	var user models.User

	if err := db.First(&user, "user_login = ?", id).Error; err != nil {
		return response.NotFound("user", id, "username")
	}

	return response.Response(user).Send(c)
}
