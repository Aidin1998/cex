package accounts

import (
	"context"
	"fmt"

	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"cex/internal/accounts/api"
	"cex/internal/accounts/db"
	"cex/pkg/cfg"
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

	// 3) Connect to the database and run migrations
	ctx := context.Background()
	database, err := db.ConnectAndMigrate(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize database: %v", err))
	}
	defer database.Close()

	// 4) Bootstrap Echo
	e := echo.New()

	// 5) Global middleware
	e.Use(middleware.Recover())      // recover panics
	e.Use(echozap.ZapLogger(zapLog)) // request/response logging
	// Removed echozap.ZapRecovery as it does not exist
	e.Use(middleware.Recover()) // recover panics

	// 6) JWT (stubbed for now—swap in your auth module when ready)
	// e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
	//     SigningKey: []byte(cfg.Cfg.JWTSecret),
	// }))

	// 7) Register account routes
	api.RegisterRoutes(e)

	// 8) Start server
	addr := fmt.Sprintf(":%s", cfg.Cfg.Accounts.Port)
	zapLog.Info("starting accounts service", zap.String("addr", addr))
	e.Logger.Fatal(e.Start(addr))
}
