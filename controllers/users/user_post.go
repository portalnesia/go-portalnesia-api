package user_controllers

import (
	"fmt"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"portalnesia.com/api/config"
	"portalnesia.com/api/models"
	"portalnesia.com/api/response"
	util "portalnesia.com/api/utils"
)

type UploadRequest struct {
	FileID string `json:"file_id" form:"file_id"`
}

func UploadPhotoProfile(c *fiber.Ctx) error {
	db := config.DB
	id := c.Params("id", "")
	ctx := c.Locals("ctx").(*models.Context)

	if ctx.User == nil || id != ctx.User.Username {
		return response.Authorization(fiber.StatusForbidden, response.ErrorAuthorizationUnauthorizedUser)
	}

	file, err := c.FormFile("image")

	if err != nil {
		var req UploadRequest
		c.BodyParser(&req)

		if req.FileID == "" {
			return response.BadParameter("file_id")
		}

		var media models.Media
		if err := db.First(&media, "unik = ?", req.FileID).Error; err != nil {
			return response.NotFound("files", req.FileID, "")
		}

		if ctx.User.ID == media.UserID && media.Jenis == "foto" && media.Sumber != nil && media.Path != nil {
			if err := db.Model(ctx.User).Where("id = ?", ctx.User.ID).Update("gambar", *media.Path).Error; err != nil {
				fmt.Println(err)
				return response.Server()
			}
			return response.Response(fiber.Map{"image": util.StaticUrl(fmt.Sprintf("img/content?image=%s", url.QueryEscape(*media.Path)))}).Send(c)
		}

		return response.InvalidParameter("file_id", "", "You cannot set profile pictures with this files")
	} else {
		filetype := file.Header.Get("Content-Type")
		if filetype != "image/jpeg" && filetype != "image/png" {
			return response.UploadError(response.ErrorUpload_FileUnsupported)
		}
		if file.Size > 5242880 {
			return response.UploadError(response.ErrorUpload_FileSize)
		}

		media, err := models.UploadImage(models.UploadImageConfig{
			File:          file,
			Context:       ctx,
			Provider:      "profile",
			Tampil:        true,
			Private:       false,
			Path:          "image",
			DynamicFolder: true,
		})

		if err != nil {
			fmt.Println(err)
			return response.Server()
		}

		if err := c.SaveFile(file, util.NewPath().Content(*media.Path)); err != nil {
			fmt.Println(err)
			return response.Server()
		}

		if err := db.Model(ctx.User).Where("id = ?", ctx.User.ID).Update("gambar", *media.Path).Error; err != nil {
			fmt.Println(err)
			return response.Server()
		}

		return response.Response(*media).Send(c)
	}
}
