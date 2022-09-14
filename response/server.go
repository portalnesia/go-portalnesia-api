package response

import "github.com/gofiber/fiber/v2"

const ErrorServer int = 500

func Server(msg ...interface{}) *Error {
	m := "internal server error"
	if len(msg) == 1 {
		switch v := msg[0].(type) {
		case string:
			if v != "" {
				m = v
			}
		}
	}
	return NewError(fiber.StatusServiceUnavailable, ErrorServer, "server", m)
}
