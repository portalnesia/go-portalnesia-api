package response

const ErrorServer int = 500

func Server(msg string) *Error {
	if msg == "" {
		msg = "internal server error"
	}
	return NewError(503, ErrorServer, "server", msg)
}
