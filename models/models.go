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

type TransactionSide string

const (
	TransactionSideDR TransactionSide = "DR"
	TransactionSideCR TransactionSide = "CR"
)

// Transaction represents the ledger entry stored in DynamoDB.
type Transaction struct {
	PK            string          `dynamodbav:"PK"`
	SK            string          `dynamodbav:"SK"`
	TransactionID string          `json:"transaction_id"  dynamodbav:"transaction_id"`
	Side          TransactionSide `json:"side" dynamodbav:"side"` // e.g., "DEBIT" or "CREDIT"
	Amount        int64           `json:"amount" dynamodbav:"amount"`
	Description   string          `json:"description" dynamodbav:"description"`
	CounterPartID *string         `json:"counter_part_id" dynamodbav:"counter_part_id"`
	AccountID     string          `json:"account_id" dynamodbav:"account_id"`
	CreatedAt     time.Time       `json:"created_at" dynamodbav:"created_at"`
}

type OutboxTransactionStatus string

const (
	OutboxTransactionStatusPENDING OutboxTransactionStatus = "PENDING"
	OutboxTransactionStatusSETTLED OutboxTransactionStatus = "SETTLED"
	OutboxTransactionStatusDEAD    OutboxTransactionStatus = "DEAD"
)

// OutboxTransaction represents a pending transaction task stored in PostgreSQL.
type OutboxTransaction struct {
	ID          string                  `json:"id"`
	Transaction `json:"transaction"`    // Embedded fields are promoted
	Status      OutboxTransactionStatus `json:"status"` // e.g., "PENDING", "COMPLETED", "FAILED"
	RetryCount  int                     `json:"retry_count"`
	LastRetryAt *time.Time              `json:"last_retry_at"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
}
