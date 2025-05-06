package integration

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"cex/internal/accounts/api"
	"cex/internal/accounts/db"
	"cex/pkg/cfg"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func TestAccountsEndToEnd(t *testing.T) {
	ctx := context.Background()
	// 1) Start CockroachDB container
	req := testcontainers.ContainerRequest{
		Image:        "cockroachdb/cockroach:v22.1.7",
		Cmd:          []string{"start-single-node", "--insecure"},
		ExposedPorts: []string{"26257/tcp"},
		WaitingFor:   wait.ForLog("node starting"),
	}
	crdb, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req, Started: true,
	})
	assert.NoError(t, err)
	defer crdb.Terminate(ctx)

	host, err := crdb.Host(ctx)
	assert.NoError(t, err)
	port, err := crdb.MappedPort(ctx, "26257")
	assert.NoError(t, err)
	dsn := fmt.Sprintf("postgresql://root@%s:%s/defaultdb?sslmode=disable", host, port.Port())

	// 2) Set cfg and run migrations
	cfg.Cfg.Accounts.DSN = dsn
	dbConn, err := db.ConnectAndMigrate(ctx, dsn)
	assert.NoError(t, err)

	// 3) Start HTTP server in background
	e := echo.New()
	e.Use(middleware.Recover())
	api.RegisterRoutes(e, dbConn)

	go e.Start(":8080")
	time.Sleep(1 * time.Second)

	// 4) Create account via HTTP
	resp, err := http.Post("http://localhost:8080/accounts", "application/json",
		strings.NewReader(`{"owner_id":"`+uuid.New().String()+`"}`))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}
