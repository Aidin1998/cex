package apiutil

import (
	"log/slog"
	"net/http"

	"cex/pkg/errors"

	"github.com/labstack/echo/v4"
)

// NotFoundError represents an error for not found resources.
type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

// BadRequestError represents an error for bad requests.
type BadRequestError struct {
	Message string
}

func (e *BadRequestError) Error() string {
	return e.Message
}

func ErrorHandler(slog *slog.Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		var statusCode *errors.StatusCode
		if errors.As(err, &statusCode) {
			if *statusCode < http.StatusInternalServerError {
				c.JSON(int(*statusCode), err)
				return
			}

			slog.ErrorContext(c.Request().Context(), err.Error())
			c.NoContent(int(*statusCode))
			return
		}

		c.NoContent(http.StatusInternalServerError)
	}
}

// JSONErrorHandler is a custom error handler for Echo that returns JSON responses.
func JSONErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.JSON(code, map[string]interface{}{
		"error": err.Error(),
	})
}

// HandleServiceError handles errors returned by the service layer
func HandleServiceError(c echo.Context, err error) error {
	// Example implementation: map service errors to HTTP responses
	switch err.(type) {
	case *BadRequestError:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	case *NotFoundError:
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	default:
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
}

// NewBadRequestError creates a new HTTP 400 Bad Request error
func NewBadRequestError(message string) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusBadRequest, message)
}
