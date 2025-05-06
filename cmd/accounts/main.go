package accounts

import (
	"fmt"

	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"cex/internal/accounts/api"
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

	// 3) Bootstrap Echo
	e := echo.New()

	// 4) Global middleware
	e.Use(middleware.Recover())      // recover panics
	e.Use(echozap.ZapLogger(zapLog)) // request/response logging
	// Removed echozap.ZapRecovery as it does not exist
	e.Use(middleware.Recover()) // recover panics

	// 5) JWT (stubbed for now—swap in your auth module when ready)
	// e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
	//     SigningKey: []byte(cfg.Cfg.JWTSecret),
	// }))

	// 6) Register account routes
	api.RegisterRoutes(e)

	// 7) Start server
	addr := fmt.Sprintf(":%s", cfg.Cfg.Accounts.Port)
	zapLog.Info("starting accounts service", zap.String("addr", addr))
	e.Logger.Fatal(e.Start(addr))
}
