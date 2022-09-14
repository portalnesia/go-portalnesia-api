package response

import "github.com/gofiber/fiber/v2"

const (
	ErrorUpload_FileSize        = 802
	ErrorUpload_FileUnsupported = 808
	ErrorUpload_Unknown         = 809
)

func getErrorMessage(err int) string {
	switch err {
	case ErrorUpload_FileSize:
		return "File too large"
	case ErrorUpload_FileUnsupported:
		return "File is not supported"
	default:
		return "Unknown upload error"
	}
}

func UploadError(errs int) error {
	msg := getErrorMessage(errs)

	return NewError(fiber.StatusBadRequest, errs, "upload", msg)
}
