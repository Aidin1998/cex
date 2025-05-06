package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Account struct {
	ID        uuid.UUID       `db:"id"`
	OwnerID   uuid.UUID       `db:"owner_id"`
	Balance   decimal.Decimal `db:"balance"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
}

// TableName is the database table for Account.
func (Account) TableName() string { return "accounts" }
