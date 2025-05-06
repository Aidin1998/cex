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

	"cex/internal/accounts/service"
)

func TestCreateAndGetAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	svc := service.NewService(db)

	ownerID := uuid.New()
	now := time.Now().UTC()

	// Expect INSERT
	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO accounts (id, owner_id, balance, created_at, updated_at)")).
		WithArgs(sqlmock.AnyArg(), ownerID, decimal.Zero, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Call CreateAccount
	acct, err := svc.CreateAccount(context.Background(), ownerID)
	assert.NoError(t, err)
	assert.Equal(t, ownerID, acct.OwnerID)
	assert.True(t, acct.Balance.Equal(decimal.Zero))

	// Expect SELECT
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, owner_id, balance, created_at, updated_at FROM accounts WHERE id = $1")).
		WithArgs(acct.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "owner_id", "balance", "created_at", "updated_at"}).
			AddRow(acct.ID, ownerID, decimal.Zero, now, now))

	// Call GetAccount
	got, err := svc.GetAccount(context.Background(), acct.ID)
	assert.NoError(t, err)
	assert.Equal(t, acct.ID, got.ID)

	// Ensure all expectations met
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateBalanceConcurrency(t *testing.T) {
	db, mock, _ := sqlmock.New()
	svc := service.NewService(db)

	id := uuid.New()
	old := decimal.NewFromInt(100)
	delta := decimal.NewFromInt(50)
	newBal := old.Add(delta)
	now := time.Now().UTC()

	// Begin tx
	mock.ExpectBegin()
	// Lock row
	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT balance FROM accounts WHERE id=$1 FOR UPDATE")).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(old))
	// Update
	mock.ExpectExec(regexp.QuoteMeta(
		"UPDATE accounts SET balance=$1, updated_at=$2 WHERE id=$3")).
		WithArgs(newBal, sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := svc.UpdateBalance(context.Background(), id, delta)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
