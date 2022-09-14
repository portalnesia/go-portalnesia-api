package response

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

const (
	ErrorBadParameter     = 700
	ErrorInvalidParameter = 710
)

func BadParameter(what_missing string) *Error {
	msg := fmt.Sprintf("Missing `%s` parameter", what_missing)

	return NewError(fiber.StatusBadRequest, ErrorBadParameter, "bad_parameter", msg)
}

func InvalidParameter(what_invalid string, should string, text string) *Error {
	msg := fmt.Sprintf("Invalid `%s` parameter", what_invalid)
	if text != "" {
		msg += fmt.Sprintf(". %s", text)
	}
	if should != "" {
		msg += fmt.Sprintf(". %s must be %s", what_invalid, should)
	}

	return NewError(fiber.StatusBadRequest, ErrorInvalidParameter, "invalid_parameter", msg)
}
