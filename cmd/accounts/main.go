package accounts

import (
	"context"
	"fmt"

	"cex/internal/accounts/api"
	"cex/internal/accounts/db"
	"cex/internal/accounts/queue"
	"cex/internal/accounts/service"
	"cex/pkg/cfg"
	"cex/pkg/otel"

	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func main() {
	// 1) Load shared config (must set cfg.Cfg.Accounts.Port first)
	cfg.Init()

	// 2) Create a Zap logger
	zapLog, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("failed to init zap: %v", err))
	}
	defer zapLog.Sync()

	// 3) Initialize OpenTelemetry tracer
	if err := otel.InitTracer("accounts"); err != nil {
		panic(fmt.Sprintf("failed to init tracer: %v", err))
	}

	// 4) Connect to the database and run migrations
	ctx := context.Background()
	database, err := db.ConnectAndMigrate(ctx, cfg.Cfg.Accounts.DSN)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize database: %v", err))
	}
	defer database.Close()

	// 5) Bootstrap Echo
	e := echo.New()

	// 6) Global middleware
	e.Use(middleware.Recover())      // recover panics
	e.Use(echozap.ZapLogger(zapLog)) // request/response logging
	// Removed echozap.ZapRecovery as it does not exist
	e.Use(middleware.Recover()) // recover panics

	// 7) JWT (stubbed for now—swap in your auth module when ready)
	// e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
	//     SigningKey: []byte(cfg.Cfg.JWTSecret),
	// }))

	// 8) Create a queue publisher
	publisher := queue.NewPublisher(cfg.Cfg.Queue.Topics, cfg.Cfg.Queue.URL)
	defer publisher.Close()

	// Create a service instance and register account routes
	serviceInstance := service.NewService(database, publisher)
	api.RegisterRoutes(e, serviceInstance)

	// 9) Start server
	addr := fmt.Sprintf(":%s", cfg.Cfg.Accounts.Port)
	zapLog.Info("starting accounts service", zap.String("addr", addr))
	e.Logger.Fatal(e.Start(addr))
}
