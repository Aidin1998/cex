package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"cex/internal/accounts/service"
	"cex/pkg/apiutil"
)

func RegisterRoutes(e *echo.Echo, svc *service.Service) {
	// 1) global middleware for JSON errors
	e.Use(middleware.Recover())
	e.HTTPErrorHandler = apiutil.JSONErrorHandler

	// 2) bind & validate
	v := validator.New()
	e.Validator = apiutil.NewEchoValidator(v)

	// 3) create a sub-router
	g := e.Group("/accounts")

	// 4) POST /accounts
	g.POST("", createAccountHandler(svc))

	// 5) GET /accounts/:id
	g.GET("/:id", getAccountHandler(svc))

	// 6) GET /accounts?owner_id=&offset=&limit=
	g.GET("", listAccountsHandler(svc))
}
