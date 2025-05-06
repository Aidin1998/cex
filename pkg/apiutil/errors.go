package apiutil

import (
	"log/slog"
	"net/http"

	"cex/pkg/errors"

	"github.com/labstack/echo/v4"
)

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
