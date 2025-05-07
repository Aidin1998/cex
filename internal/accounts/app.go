// internal/accounts/app.go
package accounts

import (
	"context"
	"database/sql"
	"net/http"

	"cex/internal/accounts/api"
	"cex/internal/accounts/db"
	"cex/internal/accounts/metrics"
	"cex/pkg/cfg"

	"github.com/brpaz/echozap"
	echoprom "github.com/labstack/echo-contrib/echoprometheus"
	jwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// NewServer bootstraps the Accounts HTTP service with logging, metrics, auth, DB, and routes.
func NewServer() (*echo.Echo, *sql.DB, error) {
	// 1) Load and validate configuration
	cfg.Init()

	// 2) Initialize Zap logger
	zapLog, err := zap.NewProduction()
	if err != nil {
		return nil, nil, err
	}

	// 3) Initialize Prometheus metrics
	metrics.InitMetrics()

	e := echo.New()
	e.Use(
		echozap.ZapLogger(zapLog),          // request logging
		middleware.Recover(),               // panic → 500 + log
		middleware.RequestID(),             // inject request IDs
		echoprom.NewMiddleware("accounts"), // Prometheus
	)

	if err := e.Start(":" + cfg.Cfg.Accounts.Port); err != nil {
		zapLog.Fatal("failed to start accounts service", zap.Error(err))
	}

	// Expose /metrics for Prometheus scraping
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	// 5) JWT authentication for all routes under /accounts
	e.Use(jwt.WithConfig(jwt.Config{
		SigningKey:  []byte(cfg.Cfg.Users.JWTSecret),
		ContextKey:  "user",                 // claims stored in c.Get("user")
		TokenLookup: "header:Authorization", // look for "Authorization: Bearer <token>"

	}))

	// 6) Connect to CockroachDB/Postgres and run migrations
	ctx := context.Background()
	dbConn, err := db.OpenAndMigrate(ctx, cfg.Cfg.Accounts.DSN)
	if err != nil {
		return nil, nil, err
	}

	// 7) Mount API routes, passing the live *sql.DB
	api.RegisterRoutes(e, dbConn)

	// 8) Health‐check endpoint
	e.GET("/healthz", func(c echo.Context) error {
		zapLog.Info("health check")
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	return e, dbConn, nil
}
