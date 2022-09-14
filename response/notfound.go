package response

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/portalnesia/go-utils"
)

const (
	ErrorNotfound         int = 401
	ErrorEndpointNotfound int = 402
)

func NotFound(tipe string, id string, idName string) *Error {
	if idName == "" {
		idName = "id"
	}

	msg := fmt.Sprintf("%s with %s `%s` not found", utils.Ucwords(tipe), idName, id)

	return NewError(fiber.StatusNotFound, ErrorNotfound, "notfound", msg)
}

func MultipleNotFound(tipe string, id []string, idName []string) *Error {
	msg := ""

	for i := 0; i < len(id); i++ {
		if i == 0 {
			msg += fmt.Sprintf("%s with %s `%s`", utils.Ucwords(tipe), idName[i], id[i])
		} else {
			msg += fmt.Sprintf(" and %s `%s`", idName[i], id[i])
		}
	}

	msg += " not found"

	return NewError(fiber.StatusNotFound, ErrorNotfound, "notfound", msg)
}

func EndpointNotFound() *Error {
	return NewError(fiber.StatusNotFound, ErrorEndpointNotfound, "notfound", "Invalid endpoint")
}
