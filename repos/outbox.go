package repos

import (
	"context"
	"fake-ledger/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OutboxRepo interface {
	Create(ctx context.Context, tx pgx.Tx, outbox *models.OutboxTransaction) (*models.OutboxTransaction, error)
	Update(ctx context.Context, outbox models.OutboxTransaction) error
	FetchPending(ctx context.Context, limit int) ([]*models.OutboxTransaction, error)
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

func (r *OutboxPostgresRepo) Update(ctx context.Context, outbox models.OutboxTransaction) error {
	query := `UPDATE outbox SET status = $1, updated_at = NOW(), retry_count = $2, last_retry_at = NOW() WHERE id = $3`
	_, err := r.pool.Exec(ctx, query, outbox.Status, outbox.RetryCount, outbox.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *OutboxPostgresRepo) FetchPending(ctx context.Context, limit int) ([]*models.OutboxTransaction, error) {
	query := `SELECT id,
				side,
				amount,
				status,
				description,
				counter_part_id,
				account_id,
				retry_count,
				created_at
			FROM outbox WHERE status = $1 LIMIT $2`
	rows, err := r.pool.Query(ctx, query, models.OutboxTransactionStatusPENDING, limit)
	if err != nil {
		return []*models.OutboxTransaction{}, err
	}
	defer rows.Close()

	var results []*models.OutboxTransaction
	for rows.Next() {
		var item models.OutboxTransaction
		err := rows.Scan(
			&item.ID,
			&item.Side,
			&item.Amount,
			&item.Status,
			&item.Description,
			&item.CounterPartID,
			&item.AccountID,
			&item.RetryCount,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, &item)
	}

	return results, nil
}
