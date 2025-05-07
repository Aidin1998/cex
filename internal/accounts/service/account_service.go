package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"cex/internal/accounts/model"
	"cex/internal/accounts/queue"
	"cex/pkg/apiutil"

	"go.opentelemetry.io/otel"
)

type AccountService struct {
	db        *sql.DB
	publisher *queue.Publisher
}

func NewAccountService(db *sql.DB, pub *queue.Publisher) *AccountService {
	return &AccountService{db: db, publisher: pub}
}

func (s *AccountService) CreateAccount(ctx context.Context, ownerID uuid.UUID, accountType string) (model.Account, error) {
	tracer := otel.Tracer("accounts-service")
	ctx, span := tracer.Start(ctx, "AccountService.CreateAccount")
	defer span.End()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return model.Account{}, err
	}

	var account model.Account
	account.ID = uuid.New()
	account.OwnerID = ownerID
	account.Balance = decimal.Zero
	account.CreatedAt = time.Now().UTC()
	account.UpdatedAt = account.CreatedAt

	_, err = tx.ExecContext(ctx, `
		INSERT INTO accounts (id, owner_id, balance, account_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		account.ID, account.OwnerID, account.Balance, accountType, account.CreatedAt, account.UpdatedAt,
	)
	if err != nil {
		tx.Rollback()
		return model.Account{}, err
	}

	ev := apiutil.AccountCreatedEvent{
		EventID:     uuid.New(),
		AccountID:   account.ID,
		OwnerID:     ownerID,
		AccountType: accountType,
		Timestamp:   time.Now().UTC(),
	}
	_ = s.publisher.PublishAccountCreated(ctx, ev)

	return account, tx.Commit()
}

func (s *AccountService) GetAccount(ctx context.Context, id uuid.UUID) (model.Account, error) {
	var account model.Account
	err := s.db.QueryRowContext(ctx, `
		SELECT id, owner_id, balance, account_type, created_at, updated_at 
		FROM accounts WHERE id = $1`, id,
	).Scan(
		&account.ID,
		&account.OwnerID,
		&account.Balance,
		&account.Type,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return model.Account{}, err
	}
	if err != nil {
		return model.Account{}, err
	}
	return account, nil
}

func (s *AccountService) ListAccounts(ctx context.Context, ownerID uuid.UUID, offset, limit int) ([]model.Account, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, owner_id, balance, account_type, created_at, updated_at 
		FROM accounts WHERE owner_id = $1 
		ORDER BY created_at DESC OFFSET $2 LIMIT $3`,
		ownerID, offset, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []model.Account
	for rows.Next() {
		var account model.Account
		if err := rows.Scan(
			&account.ID,
			&account.OwnerID,
			&account.Balance,
			&account.Type,
			&account.CreatedAt,
			&account.UpdatedAt,
		); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (s *AccountService) UpdateBalance(ctx context.Context, id uuid.UUID, delta decimal.Decimal) error {
	tracer := otel.Tracer("accounts-service")
	ctx, span := tracer.Start(ctx, "AccountService.UpdateBalance")
	defer span.End()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var oldBalance decimal.Decimal
	err = tx.QueryRowContext(ctx, `
		SELECT balance FROM accounts WHERE id = $1 FOR UPDATE`, id,
	).Scan(&oldBalance)
	if err != nil {
		tx.Rollback()
		return err
	}

	newBalance := oldBalance.Add(delta)
	_, err = tx.ExecContext(ctx, `
		UPDATE accounts SET balance = $1, updated_at = $2 WHERE id = $3`,
		newBalance, time.Now().UTC(), id,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	ev := apiutil.BalanceUpdatedEvent{
		EventID:    uuid.New(),
		AccountID:  id,
		OldBalance: oldBalance.String(),
		NewBalance: newBalance.String(),
		Delta:      delta.String(),
		Reason:     "manual-update",
		Timestamp:  time.Now().UTC(),
	}
	_ = s.publisher.PublishBalanceUpdated(ctx, ev)

	return tx.Commit()
}
