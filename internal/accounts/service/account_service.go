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

type Service struct {
	db        *sql.DB
	publisher *queue.Publisher
}

// NewService returns a new account service.
func NewService(db *sql.DB, pub *queue.Publisher) *Service {
	return &Service{db: db, publisher: pub}
}

func (s *Service) CreateAccount(ctx context.Context, ownerID uuid.UUID) (model.Account, error) {
	tracer := otel.Tracer("accounts-service")
	ctx, span := tracer.Start(ctx, "Service.CreateAccount")
	defer span.End()

	id := uuid.New()
	now := time.Now().UTC()
	initialBalance := decimal.Zero

	_, err := s.db.ExecContext(ctx, `
		INSERT INTO accounts (id, owner_id, balance, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5)`,
		id, ownerID, initialBalance, now, now,
	)
	if err != nil {
		return model.Account{}, err
	}

	account := model.Account{
		ID:        id,
		OwnerID:   ownerID,
		Balance:   initialBalance,
		CreatedAt: now,
		UpdatedAt: now,
	}

	ev := apiutil.AccountCreatedEvent{
		EventID:   uuid.New(),
		AccountID: account.ID,
		OwnerID:   account.OwnerID,
		Timestamp: time.Now().UTC(),
	}
	_ = s.publisher.PublishAccountCreated(ctx, ev)

	return account, nil
}

func (s *Service) GetAccount(ctx context.Context, id uuid.UUID) (model.Account, error) {
	var account model.Account
	err := s.db.QueryRowContext(ctx, `
		SELECT id, owner_id, balance, created_at, updated_at 
		FROM accounts WHERE id = $1`, id,
	).Scan(
		&account.ID,
		&account.OwnerID,
		&account.Balance,
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

func (s *Service) ListAccounts(ctx context.Context, ownerID uuid.UUID, offset, limit int) ([]model.Account, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, owner_id, balance, created_at, updated_at 
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
			&account.CreatedAt,
			&account.UpdatedAt,
		); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func (s *Service) UpdateBalance(ctx context.Context, id uuid.UUID, delta decimal.Decimal) error {
	tracer := otel.Tracer("accounts-service")
	ctx, span := tracer.Start(ctx, "Service.UpdateBalance")
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
