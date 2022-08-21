package response

import (
	"reflect"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type response[T any] struct {
	Error *string `json:"error"`
	Data  T       `json:"data"`
}

type paginationResponse[T any] struct {
	Page      int64 `json:"page"`
	Total     int64 `json:"total"`
	TotalPage int64 `json:"total_page"`
	CanLoad   bool  `json:"can_load"`
	Data      []T   `json:"data"`
}

func Response[T any](data T) response[T] {
	return response[T]{
		Data: data,
	}
}

func (r response[T]) Send(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(r)
}

type QueryPagination[T any] struct {
	Start   int64
	Page    int64
	PerPage int64
}

func GetPagination[T any](c *fiber.Ctx) QueryPagination[T] {
	page := c.Query("page", "1")
	perPage := c.Query("per_page", "15")
	pageInt, _ := strconv.Atoi(page)
	perPageInt, _ := strconv.Atoi(perPage)

	start := 0
	if pageInt > 1 {
		start = (pageInt * perPageInt) - perPageInt
	}

	return QueryPagination[T]{
		Start:   int64(start),
		Page:    int64(pageInt),
		PerPage: int64(perPageInt),
	}
}

func (page QueryPagination[T]) PaginationResponse(g *gorm.DB) response[paginationResponse[T]] {
	var data []T
	var total int64 = 0

	g.Limit(int(page.PerPage)).Offset(int(page.Start)).Find(&data)

	st := reflect.ValueOf(data[0])
	table := st.MethodByName("TableName").Call(nil)[0].String()

	g.Table(table).Count(&total)

	totalPage := total / page.PerPage
	canLoad := false
	if page.Page < totalPage {
		canLoad = true
	}

	return response[paginationResponse[T]]{
		Data: paginationResponse[T]{
			Page:      page.Page,
			Total:     total,
			TotalPage: totalPage,
			CanLoad:   canLoad,
			Data:      data,
		},
	}
}
