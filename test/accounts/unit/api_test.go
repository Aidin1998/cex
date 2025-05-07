package unit

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"cex/internal/accounts/api"
)

func TestRegisterRoutes(t *testing.T) {
	e := echo.New()
	dbConn := setupTestDB() // Mock or setup a test database connection
	defer dbConn.Close()

	api.RegisterRoutes(e, dbConn)

	req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assuming a handler is registered for GET /accounts
	h := func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	}
	e.GET("/accounts", h)

	if assert.NoError(t, h(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "OK", rec.Body.String())
	}
}

func setupTestDB() *sql.DB {
	// Create a mock database connection using sqlmock
	db, _, err := sqlmock.New()
	if err != nil {
		panic("failed to create mock database: " + err.Error())
	}
	return db
}
