package repos

import (
	"context"
	"fake-ledger/models"

	"github.com/jackc/pgx/v5"
)

type OutboxRepo interface {
	Create(ctx context.Context, tx pgx.Tx, outbox *models.OutboxTransaction) error
	Update(ctx context.Context, outbox models.OutboxTransaction) error
	FetchPending(ctx context.Context, limit int) ([]*models.OutboxTransaction, error)
}
