package repos

import (
	"context"
	"fake-ledger/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AccountRepo interface {
	Read(ctx context.Context, ID string) (*models.Account, error)
	Create(ctx context.Context, account *models.Account) (*models.Account, error)
	UpdateBalance(ctx context.Context, tx pgx.Tx, ID string, amount int64) error
}

type PostgresAccountRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresAccountRepo(pool *pgxpool.Pool) *PostgresAccountRepo {
	return &PostgresAccountRepo{
		pool,
	}
}

func (r *PostgresAccountRepo) Read(ctx context.Context, ID string) (*models.Account, error) {
	query := `SELECT id, account_number, nickname, balance, created_at, updated_at FROM accounts WHERE id = $1`

	row := r.pool.QueryRow(ctx, query, ID)
	account := &models.Account{}
	err := row.Scan(
		&account.ID,
		&account.AccountNumber,
		&account.Nickname,
		&account.Balance,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (r *PostgresAccountRepo) Create(ctx context.Context, account *models.Account) (*models.Account, error) {
	query := `INSERT INTO accounts (id, account_number, nickname, balance) VALUES (gen_random_uuid(), $1, $2, 0) RETURNING id, account_number, nickname, balance,created_at, updated_at `
	row := r.pool.QueryRow(ctx, query, account.AccountNumber, account.Nickname)
	created := &models.Account{}

	err := row.Scan(
		&created.ID,
		&created.AccountNumber,
		&created.Nickname,
		&created.Balance,
		&created.CreatedAt,
		&created.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return created, nil
}

func (r *PostgresAccountRepo) UpdateBalance(ctx context.Context, tx pgx.Tx, ID string, amount int64) error {
	query := `UPDATE accounts SET balance = balance + $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := tx.Exec(ctx, query, amount, ID)
	return err
}
