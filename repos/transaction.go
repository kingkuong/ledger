package repos

import (
	"context"
	"fake-ledger/models"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type TransactionRepo interface {
	FetchTransactions(ctx context.Context, accountID string, from *time.Time, to *time.Time) ([]models.Transaction, error)
	Create(ctx context.Context, transaction *models.Transaction) (*models.Transaction, error)
}

const TableName = "transactions"

type DynamoTransactionRepo struct {
	client *dynamodb.Client
}

func NewDynamoTransactionRepo(client *dynamodb.Client) *DynamoTransactionRepo {
	return &DynamoTransactionRepo{
		client: client,
	}
}

// this method needs to be idempotent
func (r *DynamoTransactionRepo) Create(ctx context.Context, transaction *models.Transaction) (*models.Transaction, error) {
	pk := generatePK(transaction.AccountID)
	sk := generateSK(&transaction.CreatedAt)

	transaction.PK = pk
	transaction.SK = sk

	item, err := attributevalue.MarshalMap(transaction)
	if err != nil {
		return &models.Transaction{}, err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item:      item,
	})

	if err != nil {
		return &models.Transaction{}, err
	}

	return transaction, nil
}

func (r *DynamoTransactionRepo) FetchTransactions(ctx context.Context, accountID string, from *time.Time, to *time.Time) ([]models.Transaction, error) {
	if accountID == "" || from == nil || to == nil {
		return []models.Transaction{}, fmt.Errorf("missing params")
	}

	pk := generatePK(accountID)
	skFrom := from.Format("2006-01-02")
	skTo := to.AddDate(0, 0, 1).Format("2006-01-02")

	// could also do `\xff`, which is the last string in UTF-8/binary
	//skTo := to.Format("2006-01-02") + `#\xff`

	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(TableName),
		KeyConditionExpression: aws.String("PK = :pk AND SK BETWEEN :from AND :to"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":   &types.AttributeValueMemberS{Value: pk},
			":from": &types.AttributeValueMemberS{Value: skFrom},
			":to":   &types.AttributeValueMemberS{Value: skTo},
		},
	})

	if err != nil {
		return []models.Transaction{}, err
	}

	var transactions []models.Transaction
	err = attributevalue.UnmarshalListOfMaps(result.Items, &transactions)
	if err != nil {
		return []models.Transaction{}, err
	}

	return transactions, nil
}

func generatePK(accountID string) string {
	return fmt.Sprintf("account#%s", accountID)
}

func generateSK(createAt *time.Time) string {
	return fmt.Sprintf("%s#%s", createAt.Format("2006-01-02"), uuid.New().String())
}
