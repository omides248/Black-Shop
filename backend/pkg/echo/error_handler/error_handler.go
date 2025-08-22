package error_handler

import (
	"errors"
	"net/http"
	"strings"
	"unicode"

	ozzo "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type ErrorResponse struct {
	StatusCode int         `json:"status_code"`
	ErrorType  string      `json:"error_type"`
	Message    string      `json:"message"`
	Details    interface{} `json:"details,omitempty"`
}

type DomainErrorMapping struct {
	StatusCode int
	Message    string
}

const (
	SystemError     = "SYSTEM_ERROR"
	ValidationError = "VALIDATION_ERROR"
	DomainError     = "DOMAIN_ERROR"
)

func NewHTTPErrorHandler(domainErrorMappings map[error]DomainErrorMapping, logger *zap.Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		var (
			statusCode = http.StatusInternalServerError
			errorType  = SystemError
			message    = "An unexpected error occurred."
			details    interface{}
		)

		// 1. Check for Validation Errors
		if errors.As(err, &ozzo.Errors{}) {
			statusCode = http.StatusBadRequest
			errorType = ValidationError
			message = "Validation failed."
			ozzoErrs := ozzo.Errors{}
			errors.As(err, &ozzoErrs)

			mapped := make(map[string]string)
			for field, e := range ozzoErrs {
				jsonKey := toSnakeCase(field)
				mapped[jsonKey] = e.Error()
			}
			details = mapped

		} else {
			// 2. Check for Domain Errors using the provided map
			for domainErr, mapping := range domainErrorMappings {
				if errors.Is(err, domainErr) {
					statusCode = mapping.StatusCode
					errorType = DomainError
					message = mapping.Message
					break
				}
			}
		}

		// 3. Fallback for other Echo or system errors
		if statusCode == http.StatusInternalServerError && c.Response().Committed {
			logger.Error("response already committed", zap.Error(err))
			return
		}

		logger.Error("handling http error",
			zap.Error(err),
			zap.Int("status_code", statusCode),
			zap.String("error_type", errorType),
			zap.Any("details", details),
		)

		_ = c.JSON(statusCode, ErrorResponse{
			StatusCode: statusCode,
			ErrorType:  errorType,
			Message:    message,
			Details:    details,
		})
	}
}

func toSnakeCase(s string) string {
	var sb strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				sb.WriteRune('_')
			}
			sb.WriteRune(unicode.ToLower(r))
		} else {
			sb.WriteRune(r)
		}
	}

	return sb.String()
}
