package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Account struct {
	ID        uuid.UUID       `db:"id" json:"id"`
	OwnerID   uuid.UUID       `db:"owner_id" json:"owner_id"`
	Type      string          `db:"type" json:"type"`
	Balance   decimal.Decimal `db:"balance" json:"balance"`
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt time.Time       `db:"updated_at" json:"updated_at"`
}

// TableName is the database table for Account.
func (Account) TableName() string { return "accounts" }
