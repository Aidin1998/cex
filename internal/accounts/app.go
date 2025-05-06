// internal/accounts/app.go
package accounts

import (
	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"cex/internal/accounts/api"
)

// NewServer returns an Echo configured with middleware & routes.
// Configuration & zap-logger are handled in cmd/accounts/main.go.
func NewServer(zapLog *zap.Logger) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(echozap.ZapLogger(zapLog))
	e.Use(middleware.Recover())
	api.RegisterRoutes(e)
	return e
}
