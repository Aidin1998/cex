package api

import (
	"net/http"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"

	"cex/internal/accounts/service"
	"cex/pkg/apiutil"
)

// createAccountRequest defines POST body
type createAccountRequest struct {
	OwnerID uuid.UUID `json:"owner_id" validate:"required"`
}

// accountResponse defines JSON output
type accountResponse struct {
	ID        uuid.UUID       `json:"id"`
	OwnerID   uuid.UUID       `json:"owner_id"`
	Balance   decimal.Decimal `json:"balance"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
}

func createAccountHandler(svc *service.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 1) bind & validate
		var req createAccountRequest
		if err := apiutil.BindAndValidate(c, &req); err != nil {
			return err
		}
		// 2) call service
		acct, err := svc.CreateAccount(c.Request().Context(), req.OwnerID)
		if err != nil {
			return apiutil.HandleServiceError(c, err)
		}
		// 3) build response
		res := accountResponse{
			ID:        acct.ID,
			OwnerID:   acct.OwnerID,
			Balance:   acct.Balance,
			CreatedAt: acct.CreatedAt.Format(time.RFC3339),
			UpdatedAt: acct.UpdatedAt.Format(time.RFC3339),
		}
		// 4) set Location header + 201
		c.Response().Header().Set("Location", "/accounts/"+acct.ID.String())
		return c.JSON(http.StatusCreated, res)
	}
}

func getAccountHandler(svc *service.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Extract userID from JWT claims
		token := c.Get("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		uidStr := claims["sub"].(string) // assume sub = user UUID
		userID, err := uuid.Parse(uidStr)
		if err != nil {
			return apiutil.NewUnauthorizedError("invalid JWT subject")
		}

		// 1) parse & validate path param
		idStr := c.Param("id")
		acctID, err := uuid.Parse(idStr)
		if err != nil {
			return apiutil.NewBadRequestError("invalid account ID")
		}

		// 2) call service
		acct, err := svc.GetAccount(c.Request().Context(), acctID)
		if err != nil {
			return apiutil.HandleServiceError(c, err)
		}

		// Check if the account belongs to the user
		if acct.OwnerID != userID {
			return apiutil.NewForbiddenError("not your account")
		}

		// 3) respond
		res := accountResponse{
			ID:        acct.ID,
			OwnerID:   acct.OwnerID,
			Balance:   acct.Balance,
			CreatedAt: acct.CreatedAt.Format(time.RFC3339),
			UpdatedAt: acct.UpdatedAt.Format(time.RFC3339),
		}
		return c.JSON(http.StatusOK, res)
	}
}

func listAccountsHandler(svc *service.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Extract userID from JWT claims
		token := c.Get("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		uidStr := claims["sub"].(string)
		userID, err := uuid.Parse(uidStr)
		if err != nil {
			return apiutil.NewUnauthorizedError("invalid JWT subject")
		}

		// 1) parse query params offset, limit
		offset, _ := strconv.Atoi(c.QueryParam("offset"))
		limit, _ := strconv.Atoi(c.QueryParam("limit"))

		// 2) call service with the userID
		accts, err := svc.ListAccounts(c.Request().Context(), userID, offset, limit)
		if err != nil {
			return apiutil.HandleServiceError(c, err)
		}

		// 3) map to []accountResponse
		var out []accountResponse
		for _, a := range accts {
			out = append(out, accountResponse{
				ID:        a.ID,
				OwnerID:   a.OwnerID,
				Balance:   a.Balance,
				CreatedAt: a.CreatedAt.Format(time.RFC3339),
				UpdatedAt: a.UpdatedAt.Format(time.RFC3339),
			})
		}
		return c.JSON(http.StatusOK, out)
	}
}
