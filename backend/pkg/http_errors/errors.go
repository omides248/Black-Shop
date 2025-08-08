package http_errors

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ErrorResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

func HandleValidationError(c echo.Context, err error) error {
	var validationErrors validation.Errors
	if errors.As(err, &validationErrors) {

		errResponse := ErrorResponse{
			Code:    "VALIDATION_ERROR",
			Message: "The given data was invalid.",
			Errors:  validationErrors,
		}
		return c.JSON(http.StatusUnprocessableEntity, errResponse)
	}

	errResponse := ErrorResponse{
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "An unexpected error occurred.",
	}
	return c.JSON(http.StatusInternalServerError, errResponse)
}
