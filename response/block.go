package response

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const ErrorBlock = 300

func Block(tipe string, id string, idName string) *Error {
	if idName == "" {
		idName = "id"
	}

	msg := fmt.Sprintf("Your %s with %s `%s` is blocked", strings.ToLower(tipe), idName, id)

	return NewError(fiber.StatusNotFound, ErrorBlock, "block", msg)
}
