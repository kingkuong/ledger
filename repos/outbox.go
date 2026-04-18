package repos

import (
	"context"
	"fake-ledger/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OutboxRepo interface {
	Create(ctx context.Context, tx pgx.Tx, outbox *models.OutboxTransaction) (*models.OutboxTransaction, error)
	//Update(ctx context.Context, outbox models.OutboxTransaction) error
	//FetchPending(ctx context.Context, limit int) ([]*models.OutboxTransaction, error)
}

type OutboxPostgresRepo struct {
	pool *pgxpool.Pool
}

func NewOutboxPostgresRepo(pool *pgxpool.Pool) *OutboxPostgresRepo {
	return &OutboxPostgresRepo{
		pool: pool,
	}
}

func (r *OutboxPostgresRepo) Create(ctx context.Context, tx pgx.Tx, outbox *models.OutboxTransaction) (*models.OutboxTransaction, error) {
	query := `INSERT INTO outbox(
		id,
		side,
		amount,
		description,
		counter_part_id,
		account_id,
		status,
		created_at,
		retry_count
	) VALUES
	(
		$1, $2, $3, $4, $5, $6, $7, now(), 0
	) RETURNING created_at, retry_count;`

	err := tx.QueryRow(ctx, query,
		outbox.ID,
		outbox.Side,
		outbox.Amount,
		outbox.Description,
		outbox.CounterPartID,
		outbox.AccountID,
		outbox.Status,
	).Scan(
		&outbox.CreatedAt,
		&outbox.RetryCount,
	)

	if err != nil {
		return nil, err
	}

	return outbox, nil
}
