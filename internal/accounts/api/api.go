package api

import (
	"database/sql"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"cex/internal/accounts/api/handlers"
	"cex/internal/accounts/service"
	"cex/pkg/apiutil"
)

var jwtMiddleware = jwt.WithConfig(jwt.Config{
	SigningKey:  []byte("your-jwt-secret"),
	ContextKey:  "user",
	TokenLookup: "header:Authorization",
})

func RegisterRoutes(e *echo.Echo, db *sql.DB) {
	// 1) global middleware for JSON errors
	e.Use(middleware.Recover())
	e.HTTPErrorHandler = apiutil.JSONErrorHandler

	// 2) bind & validate
	v := validator.New()
	e.Validator = apiutil.NewEchoValidator(v)

	// 3) create a sub-router
	// Create a queue.Publisher instance (replace with actual implementation)
	publisher := service.NewQueuePublisher() // Ensure NewQueuePublisher is implemented in the service package

	// Pass the publisher to the service
	svc := service.NewAccountService(db, publisher)
	g := e.Group("/accounts", jwtMiddleware)

	// 4) POST /accounts
	g.POST("", handlers.CreateAccountHandler(svc))

	// 5) GET /accounts/:id
	g.GET("/:id", handlers.GetAccountHandler(svc))

	// 6) GET /accounts?owner_id=&offset=&limit=
	g.GET("", handlers.ListAccountsHandler(svc))
}
