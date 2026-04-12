package service

import (
	"context"
	"fake-ledger/models"
	"fake-ledger/repos"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LedgerService struct {
	accountRepo     repos.AccountRepo
	transactionRepo repos.TransactionRepo
}

func NewLedgerService(ctx context.Context) (*LedgerService, error) {
	pool, err := pgxpool.New(ctx, os.Getenv("APP_POSTGRES_URL"))
	if err != nil {
		return nil, err
	}
	accountRepo := repos.NewPostgresAccountRepo(pool)

	client, err := NewDynamoClient(ctx, "http://dynamo:8000")
	if err != nil {
		return nil, err
	}
	transactionRepo := repos.NewDynamoTransactionRepo(client)

	return &LedgerService{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
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
	return s.transactionRepo.Create(ctx, transaction)
}
