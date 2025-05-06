package apiutil

import (
	"time"

	"github.com/google/uuid"
)

// AccountCreatedEvent is published when a new Account is created.
type AccountCreatedEvent struct {
	EventID   uuid.UUID `json:"event_id"`
	AccountID uuid.UUID `json:"account_id"`
	OwnerID   uuid.UUID `json:"owner_id"`
	Timestamp time.Time `json:"timestamp"`
}

// BalanceUpdatedEvent is published when an Account balance changes.
type BalanceUpdatedEvent struct {
	EventID    uuid.UUID `json:"event_id"`
	AccountID  uuid.UUID `json:"account_id"`
	OldBalance string    `json:"old_balance"` // decimal as string
	NewBalance string    `json:"new_balance"`
	Delta      string    `json:"delta"`
	Reason     string    `json:"reason"` // e.g. "credit", "debit"
	Timestamp  time.Time `json:"timestamp"`
}
