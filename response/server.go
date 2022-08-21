package response

import "github.com/gofiber/fiber/v2"

const ErrorServer int = 500

func Server(msg string) *Error {
	if msg == "" {
		msg = "internal server error"
	}
	return NewError(fiber.StatusServiceUnavailable, ErrorServer, "server", msg)
}
