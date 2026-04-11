package repos

import (
	"context"
	"fake-ledger/models"
	"time"
)

type TransactionRepo interface {
	FetchTransactions(ctx context.Context, accountID string, from *time.Time, to *time.Time) ([]models.Transaction, error)
	Create(ctx context.Context, transaction models.Transaction) (models.Transaction, error)
}
