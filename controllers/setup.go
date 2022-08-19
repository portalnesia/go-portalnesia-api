package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetPagination(c *fiber.Ctx) (int, int, int) {
	page := c.Query("page", "1")
	perPage := c.Query("per_page", "15")
	pageInt, _ := strconv.Atoi(page)
	perPageInt, _ := strconv.Atoi(perPage)

	start := 0
	if pageInt > 1 {
		start = (pageInt * perPageInt) - perPageInt
	}

	return start, perPageInt, pageInt
}

//start = page > 1 ? (page*per_page)-per_page : 0;
