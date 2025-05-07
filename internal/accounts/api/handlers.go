package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shopspring/decimal"

	"cex/internal/accounts/metrics"
	"cex/internal/accounts/service"
	"cex/pkg/apiutil"
)

var validate = validator.New()

// accountResponse defines JSON output
type accountResponse struct {
	ID        uuid.UUID       `json:"id"`
	OwnerID   uuid.UUID       `json:"owner_id"`
	Balance   decimal.Decimal `json:"balance"`
	Type      string          `json:"type"`
	CreatedAt string          `json:"created_at"`
	UpdatedAt string          `json:"updated_at"`
}

func CreateAccountHandler(svc service.AccountService) echo.HandlerFunc {
	type req struct {
		Type string `json:"type" validate:"required,oneof=fiat spot futures"`
	}
	return func(c echo.Context) error {
		var r req
		if err := c.Bind(&r); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
		}
		if err := validate.Struct(&r); err != nil {
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{"error": err.Error()})
		}
		userID := c.Get("userID").(string)
		userUUID, err := uuid.Parse(userID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid user ID"})
		}
		acct, err := svc.CreateAccount(c.Request().Context(), userUUID, r.Type)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
		return c.JSON(http.StatusCreated, acct)
	}
}

func getAccountHandler(svc *service.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		timer := prometheus.NewTimer(metrics.RequestDuration.WithLabelValues(c.Request().Method, c.Path()))
		defer timer.ObserveDuration()

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
			Type:      acct.Type,
			CreatedAt: acct.CreatedAt.Format(time.RFC3339),
			UpdatedAt: acct.UpdatedAt.Format(time.RFC3339),
		}

		metrics.RequestsTotal.WithLabelValues(c.Request().Method, c.Path(), strconv.Itoa(c.Response().Status)).Inc()
		return c.JSON(http.StatusOK, res)
	}
}

func listAccountsHandler(svc *service.Service) echo.HandlerFunc {
	return func(c echo.Context) error {
		timer := prometheus.NewTimer(metrics.RequestDuration.WithLabelValues(c.Request().Method, c.Path()))
		defer timer.ObserveDuration()

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
				Type:      a.Type,
				CreatedAt: a.CreatedAt.Format(time.RFC3339),
				UpdatedAt: a.UpdatedAt.Format(time.RFC3339),
			})
		}

		metrics.RequestsTotal.WithLabelValues(c.Request().Method, c.Path(), strconv.Itoa(c.Response().Status)).Inc()
		return c.JSON(http.StatusOK, out)
	}
}
