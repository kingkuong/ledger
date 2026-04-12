package repos

import (
	"context"
	"fake-ledger/models"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
)

type TransactionRepo interface {
	FetchTransactions(ctx context.Context, accountID string, from *time.Time, to *time.Time) ([]models.Transaction, error)
	Create(ctx context.Context, transaction *models.Transaction) (*models.Transaction, error)
}

type DynamoTransactionRepo struct {
	client *dynamodb.Client
}

func NewDynamoTransactionRepo(client *dynamodb.Client) *DynamoTransactionRepo {
	return &DynamoTransactionRepo{
		client: client,
	}
}

func (r *DynamoTransactionRepo) FetchTransactions(ctx context.Context, accountID string, from *time.Time, to *time.Time) ([]models.Transaction, error) {
	return []models.Transaction{}, nil
}

func (r *DynamoTransactionRepo) Create(ctx context.Context, transaction *models.Transaction) (*models.Transaction, error) {
	pk := generatePK(transaction.AccountID)
	now := time.Now()
	sk := generateSK(&now)

	transaction.PK = pk
	transaction.SK = sk
	transaction.CreatedAt = now

	item, err := attributevalue.MarshalMap(transaction)
	if err != nil {
		return &models.Transaction{}, err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("transactions"),
		Item:      item,
	})

	if err != nil {
		return &models.Transaction{}, err
	}

	return transaction, nil
}

func generatePK(accountID string) string {
	return fmt.Sprintf("account#%s", accountID)
}

func generateSK(createAt *time.Time) string {
	return fmt.Sprintf("%s#%s", createAt.Format("2006-01-02"), uuid.New().String())
}
