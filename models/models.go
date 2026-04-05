package models

import "time"

// Account represents a user's account stored in PostgreSQL.
type Account struct {
	ID            string    `json:"id"`
	AccountNumber string    `json:"account_number"`
	Nickname      *string   `json:"nickname"`
	Balance       int64     `json:"balance"` // Stored in smallest unit (e.g., cents)
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// Transaction represents the ledger entry stored in DynamoDB.
type Transaction struct {
	TransactionID string    `json:"transaction_id"`
	Side          string    `json:"side"` // e.g., "DEBIT" or "CREDIT"
	Amount        int64     `json:"amount"`
	Description   string    `json:"description"`
	CounterPartID string    `json:"counter_part_id"`
	AccountID     string    `json:"account_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// OutboxTransaction represents a pending transaction task stored in PostgreSQL.
// By embedding the Transaction struct, we ensure all data required to create 
// the DynamoDB record is captured atomically within the same Postgres transaction 
// as the Account balance update.
type OutboxTransaction struct {
	ID          string     `json:"id"`
	Transaction          `json:"transaction"` // Embedded fields are promoted
	Status      string     `json:"status"`     // e.g., "PENDING", "COMPLETED", "FAILED"
	RetryCount  int        `json:"retry_count"`
	LastRetryAt *time.Time `json:"last_retry_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
