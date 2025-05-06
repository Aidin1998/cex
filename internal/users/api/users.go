package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HelloWorldResponse represents the response for the hello world endpoint
type HelloWorldResponse struct {
	Message string `json:"message"`
}

// HelloWorld handles GET /hello endpoint
// @Summary Hello world endpoint
// @Description Returns a simple hello world message
// @Tags public
// @Accept json
// @Produce json
// @Success 200 {object} HelloWorldResponse
// @Router /hello [get]
func (a *API) HelloWorld(c echo.Context) error {
	return c.JSON(http.StatusOK, HelloWorldResponse{
		Message: "Hello, World!",
	})
}

