package accounts

import (
	"net/http"

	_ "github.com/lib/pq" // Replace with your database driver

	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"cex/internal/accounts/api"
	"cex/pkg/cfg"
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
		middleware.Recover(),
		middleware.RequestID(),
	)
	// Request‐logging + panic‐recovery via echozap
	e.Use(
		echozap.ZapLogger(zapLog),
		middleware.Recover(),
	)

	// 4) Connect to CockroachDB (Postgres) via your dbpkg helper
	dbConn, err := db.NewDB(cfg.Cfg.Accounts.DSN)
	if err != nil {
		return nil, err
	}
	// 5) Run migrations (000001_create_accounts_table.sql, etc.)
	if err := dbpkg.Migrate(dbConn); err != nil {
		return nil, err
	}

	// 6) Mount your API routes, passing in the live *sql.DB
	api.RegisterRoutes(e, dbConn)

	// 7) Health‐check endpoint
	e.GET("/healthz", func(c echo.Context) error {
		zapLog.Info("health check")
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	return e, nil
}
