package response

type Error struct {
	// HTTP Status
	Status      int    `json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Code        int    `json:"code"`
}

func (e *Error) Error() string {
	return e.Description
}

/*
const ErrorRatelimit = 200
const ErrorBlock = 300
const ErrorForbidden = 500
const ErrorBadParameter = 700
const ErrorInvalidParameter = 710
const ErrorUploadError = 800
const ErrorCustom = 900
const ErrorUnsplash = 910
const ErrorPixabay = 920
*/

func NewError(status int, code int, message string, description string) *Error {
	return &Error{
		Status:      status,
		Code:        code,
		Name:        message,
		Description: description,
	}
}
