package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes wires up all /accounts endpoints.
func RegisterRoutes(e *echo.Echo) {
	g := e.Group("/accounts")

	g.GET("/healthz", healthz)
	g.GET("", ListAccounts)
	g.POST("", CreateAccount)
	// TODO: more endpoints here...
}

func healthz(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
}

func ListAccounts(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"stub": "ListAccounts"})
}

func CreateAccount(c echo.Context) error {
	return c.JSON(http.StatusCreated, map[string]string{"stub": "CreateAccount"})
}
