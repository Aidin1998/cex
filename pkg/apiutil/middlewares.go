package apiutil

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// EchoValidator is a custom validator for Echo framework
type EchoValidator struct {
	validator *validator.Validate
}

// NewEchoValidator creates a new EchoValidator
func NewEchoValidator(v *validator.Validate) *EchoValidator {
	return &EchoValidator{validator: v}
}

// Validate validates the input struct
func (ev *EchoValidator) Validate(i interface{}) error {
	return ev.validator.Struct(i)
}

// BindAndValidate binds the request body to the given struct and validates it.
func BindAndValidate(c echo.Context, v interface{}) error {
	if err := c.Bind(v); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").SetInternal(err)
	}
	if err := validator.New().Struct(v); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "validation failed").SetInternal(err)
	}
	return nil
}
