package accounts

import (
	"context"
	"net/http"

	_ "github.com/lib/pq" // Replace with your database driver

	"cex/internal/accounts/api"
	"cex/internal/accounts/db"
	"cex/internal/accounts/metrics"
	"cex/internal/accounts/queue"
	"cex/internal/accounts/service"
	"cex/pkg/cfg"

	"github.com/brpaz/echozap"
	"github.com/labstack/echo-contrib/echoprometheus"
	jwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func NewServer() (*echo.Echo, error) {
	// 1) Load shared config (must set cfg.Cfg.Accounts.Port!)
	cfg.Init()

	// 2) Create a real Zap logger
	zapLog, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	// 3) Create Echo and wire up middleware
	e := echo.New()
	e.Use(
		echozap.ZapLogger(zapLog),
		middleware.Recover(),
		middleware.RequestID(),
	)

	// Initialize and register Prometheus metrics
	metrics.InitMetrics()

	// Prometheus middleware to track requests & latencies
	e.Use(echoprometheus.NewMiddleware("accounts"))

	// 4) JWT authentication for all /accounts routes
	e.Use(jwt.WithConfig(jwt.Config{
		SigningKey:  []byte(cfg.Cfg.Users.JWTSecret),
		ContextKey:  "user", // JWT claims will be available under c.Get("user")
		TokenLookup: "header:Authorization",
	}))

	// 5) Run migrations (000001_create_accounts_table.sql, etc.)
	ctx := context.Background()
	dbConn, err := db.OpenAndMigrate(ctx, cfg.Cfg.Accounts.DSN)
	if err != nil {
		return nil, err
	}

	// 6) Mount your API routes, passing in the live *service.Service
	publisher := queue.NewPublisher([]string{"queue1", "queue2"}, "defaultQueue") // Create a new queue.Publisher instance
	service := service.NewService(dbConn, publisher)                              // Create a new service instance
	api.RegisterRoutes(e, service)

	// 7) Health‐check endpoint
	e.GET("/healthz", func(c echo.Context) error {
		zapLog.Info("health check")
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	// Expose /metrics for Prometheus scraping
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	return e, nil
}
