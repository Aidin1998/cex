package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all accounts-related endpoints on the provided group.

// GetAccount returns one account by ID (placeholder).
func GetAccount(c echo.Context) error {
	id := c.Param("account_id")
	// TODO: implement fetching by ID
	return c.JSON(http.StatusOK, map[string]string{"stub": "GetAccount", "id": id})
}

// TransferFunds moves funds between fiat/spot/futures (placeholder).
func TransferFunds(c echo.Context) error {
	// TODO: parse JSON body, call service layer
	return c.JSON(http.StatusAccepted, map[string]string{"stub": "TransferFunds"})
}
