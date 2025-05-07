package unit

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"cex/internal/accounts/queue"
	"cex/internal/accounts/service"
)

func TestCreateAndGetAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	publisher := queue.NewPublisher([]string{"localhost:9092"}, "accounts-events")
	svc := service.NewAccountService(db, publisher)

	ownerID := uuid.New()
	now := time.Now().UTC()

	// Mock transaction begin
	mock.ExpectBegin()

	// Expect INSERT
	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO accounts (id, owner_id, balance, account_type, created_at, updated_at)")).
		WithArgs(sqlmock.AnyArg(), ownerID, decimal.Zero, "fiat", now, now).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock transaction commit
	mock.ExpectCommit()

	// Call CreateAccount
	acct, err := svc.CreateAccount(context.Background(), ownerID, "fiat")
	assert.NoError(t, err)
	assert.Equal(t, ownerID, acct.OwnerID)
	assert.True(t, acct.Balance.Equal(decimal.Zero))

	// Expect SELECT
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, owner_id, balance, account_type, created_at, updated_at FROM accounts WHERE id = $1")).
		WithArgs(acct.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "owner_id", "balance", "account_type", "created_at", "updated_at"}).
			AddRow(acct.ID, ownerID, decimal.Zero, "fiat", now, now))

	// Call GetAccount
	got, err := svc.GetAccount(context.Background(), acct.ID)
	assert.NoError(t, err)
	assert.Equal(t, acct.ID, got.ID)

	// Ensure all expectations met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateBalanceConcurrency(t *testing.T) {
	db, mock, _ := sqlmock.New()
	publisher := queue.NewPublisher([]string{"localhost:9092"}, "accounts-events")
	svc := service.NewAccountService(db, publisher)

	id := uuid.New()
	old := decimal.NewFromInt(100)
	delta := decimal.NewFromInt(50)
	newBal := old.Add(delta)

	// Begin tx
	mock.ExpectBegin()
	// Lock row
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT balance FROM accounts WHERE id = $1 FOR UPDATE")).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(old))
	// Update
	mock.ExpectExec(regexp.QuoteMeta(
		"UPDATE accounts SET balance = $1, updated_at = $2 WHERE id = $3")).
		WithArgs(newBal, sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := svc.UpdateBalance(context.Background(), id, delta)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateAccount_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	publisher := queue.NewPublisher([]string{"localhost:9092"}, "accounts-events")
	svc := service.NewAccountService(db, publisher)
	ownerID := uuid.New()
	now := time.Now().UTC()

	// Mock transaction begin
	mock.ExpectBegin()

	// Expect INSERT ... RETURNING ...
	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO accounts (id, owner_id, balance, account_type, created_at, updated_at)")).
		WithArgs(sqlmock.AnyArg(), ownerID, decimal.Zero, "fiat", now, now).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock transaction commit
	mock.ExpectCommit()

	account, err := svc.CreateAccount(context.Background(), ownerID, "fiat")
	require.NoError(t, err)
	require.Equal(t, ownerID, account.OwnerID)

	require.NoError(t, mock.ExpectationsWereMet())
}
