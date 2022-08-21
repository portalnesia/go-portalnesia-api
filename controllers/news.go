package controllers

import (
	"net/url"

	"portalnesia.com/api/config"
	"portalnesia.com/api/models"
	"portalnesia.com/api/response"

	"github.com/gofiber/fiber/v2"
)

func FindNews(c *fiber.Ctx) error {
	db := config.DB
	g := db.Order("id desc")
	return response.GetPagination[models.NewsPagination](c).PaginationResponse(g).Send(c)
}

func FindOneNews(c *fiber.Ctx) error {
	db := config.DB
	source := c.Params("source", "")
	title := c.Params("title", "")
	title, _ = url.QueryUnescape(title)

	var news models.News

	if err := db.First(&news, "source = ? AND title = ?", source, title).Error; err != nil {
		return response.MultipleNotFound("news", []string{source, title}, []string{"source", "title"})
	}
	return response.Response(news).Send(c)
}
