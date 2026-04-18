package service

import (
	"context"
	"fake-ledger/models"
	"fake-ledger/repos"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LedgerService struct {
	accountRepo           repos.AccountRepo
	transactionRepo       repos.TransactionRepo
	outboxTransactionRepo repos.OutboxRepo
	pool                  *pgxpool.Pool
}

func NewLedgerService(ctx context.Context,
	accountRepo repos.AccountRepo,
	transactionRepo repos.TransactionRepo,
	outboxTransactionRepo repos.OutboxRepo,
	pool *pgxpool.Pool,
) (*LedgerService, error) {
	return &LedgerService{
		accountRepo:           accountRepo,
		transactionRepo:       transactionRepo,
		outboxTransactionRepo: outboxTransactionRepo,
		pool:                  pool,
	}, nil
}

func NewDynamoClient(ctx context.Context, endpoint string) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(os.Getenv("AWS_DEFAULT_REGION")),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
				SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
			},
		}),
	)

	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		if endpoint != "" {
			o.BaseEndpoint = &endpoint
		}
	})

	return client, nil
}

func (s *LedgerService) GetAccount(ctx context.Context, ID string) (*models.Account, error) {
	return s.accountRepo.Read(ctx, ID)
}

func (s *LedgerService) CreateAccount(ctx context.Context, account *models.Account) (*models.Account, error) {
	return s.accountRepo.Create(ctx, account)
}

func (s *LedgerService) CreateTransaction(ctx context.Context, transaction *models.Transaction) (*models.Transaction, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	account, err := s.accountRepo.ReadForUpdate(ctx, tx, transaction.AccountID)
	if err != nil {
		return nil, err
	}

	balance := account.Balance
	amount := transaction.Amount
	if transaction.Side == models.TransactionSideDR {
		if balance < transaction.Amount {
			return nil, fmt.Errorf("Unable to commit transaction, Not Enough Balance")
		}
		amount = -amount
	}
	err = s.accountRepo.UpdateBalance(ctx, tx, transaction.AccountID, amount)
	if err != nil {
		return nil, err
	}

	_, err = s.outboxTransactionRepo.Create(ctx, tx, &models.OutboxTransaction{
		ID:          uuid.NewString(),
		Transaction: *transaction,
		Status:      models.OutboxTransactionStatusPENDING,
	})
	if err != nil {
		return nil, err
	}

	tx.Commit(ctx)

	return transaction, nil
}

func (s *LedgerService) FetchTransactions(ctx context.Context, accountID string, from *time.Time, to *time.Time) ([]models.Transaction, error) {
	return s.transactionRepo.FetchTransactions(ctx, accountID, from, to)
}

func (s *LedgerService) FetchOutboxPending(ctx context.Context, limit int) ([]*models.OutboxTransaction, error) {
	return s.outboxTransactionRepo.FetchPending(ctx, limit)
}

func (s *LedgerService) Transfer(ctx context.Context, accountIDfrom string, accountIDto string, amount int64) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	accountFrom, err := s.accountRepo.ReadForUpdate(ctx, tx, accountIDfrom)
	if err != nil {
		return err
	}
	accountTo, err := s.accountRepo.ReadForUpdate(ctx, tx, accountIDto)
	if err != nil {
		return err
	}

	if accountFrom.Balance < amount {
		return fmt.Errorf("Not enough balance")
	}

	err = s.accountRepo.UpdateBalance(ctx, tx, accountFrom.ID, -amount)
	if err != nil {
		return err
	}
	err = s.accountRepo.UpdateBalance(ctx, tx, accountTo.ID, amount)
	if err != nil {
		return err
	}
	// create outbox transactions
	// Decided to create UUID in app as we'll save trip to db for writing
	// CR transaction

	crTransactionID := uuid.NewString()
	drTransactionID := uuid.NewString()
	crTransactionOubox := &models.OutboxTransaction{
		ID: crTransactionID,
		Transaction: models.Transaction{
			Side:          models.TransactionSideCR,
			Amount:        amount,
			Description:   fmt.Sprintf("Transfer from %s", accountIDfrom),
			AccountID:     accountIDto,
			CounterPartID: &drTransactionID,
		},
		Status: models.OutboxTransactionStatusPENDING,
	}
	crTransactionOubox, err = s.outboxTransactionRepo.Create(ctx, tx, crTransactionOubox)
	if err != nil {
		return err
	}
	// DR transaction
	drTransactionOubox := &models.OutboxTransaction{
		ID: drTransactionID,
		Transaction: models.Transaction{
			Side:          models.TransactionSideDR,
			Amount:        amount,
			Description:   fmt.Sprintf("Transfer to %s", accountIDto),
			AccountID:     accountIDfrom,
			CounterPartID: &crTransactionID,
		},
		Status: models.OutboxTransactionStatusPENDING,
	}
	drTransactionOubox, err = s.outboxTransactionRepo.Create(ctx, tx, drTransactionOubox) // DR transaction
	if err != nil {
		return err
	}
	// Update account_

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
