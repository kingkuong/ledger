package service

import (
	"context"
	"fake-ledger/models"
	"fake-ledger/repos"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LedgerService struct {
	accountRepo repos.AccountRepo
}

func NewLedgerService(ctx context.Context) (*LedgerService, error) {
	pool, err := pgxpool.New(ctx, os.Getenv("APP_POSTGRES_URL"))
	if err != nil {
		return nil, err
	}

	accountRepo := repos.NewPostgresAccountRepo(pool)
	return &LedgerService{
		accountRepo,
	}, nil
}

func (s *LedgerService) GetAccount(ctx context.Context, ID string) (*models.Account, error) {
	return s.accountRepo.Read(ctx, ID)
}

func (s *LedgerService) CreateAccount(ctx context.Context, account *models.Account) (*models.Account, error) {
	return s.accountRepo.Create(ctx, account)
}
