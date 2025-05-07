package api

import (
	"cex/internal/accounts/queue"
	"database/sql"

	jwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/go-playground/validator/v10"

	"cex/internal/accounts/service"
	"cex/pkg/apiutil"
	"cex/pkg/cfg"
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

	// Create a Kafka publisher instance
	publisher := queue.NewPublisher(
		cfg.Cfg.Kafka.Brokers,       // e.g. []string{"localhost:9092"}
		cfg.Cfg.Kafka.TopicAccounts, // e.g. "accounts-events"
	)
	// Pass the publisher to the service
	svc := service.NewAccountService(db, publisher)
	// apply JWT middleware (imported above)
	g := e.Group("/accounts", jwt.WithConfig(jwt.Config{
		SigningKey:  []byte("your-jwt-secret"),
		ContextKey:  "user",
		TokenLookup: "header:Authorization",
	}))
	// 4) POST /accounts
	g.POST("", CreateAccountHandler(svc))

	// 5) GET /accounts/:id
	g.GET("/:id", GetAccountHandler(svc))

	// 6) GET /accounts?owner_id=&offset=&limit=
	g.GET("", ListAccountsHandler(svc))
}
