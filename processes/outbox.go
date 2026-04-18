package processes

import (
	"context"
	"fake-ledger/models"
	"fake-ledger/repos"
	"fmt"
	"time"
)

type OutboxProcess struct {
	outboxRepo      repos.OutboxPostgresRepo
	transactionRepo repos.DynamoTransactionRepo
}

func NewOutboxProcess(
	outboxRepo repos.OutboxPostgresRepo,
	transactionRepo repos.DynamoTransactionRepo,
) *OutboxProcess {
	return &OutboxProcess{
		outboxRepo:      outboxRepo,
		transactionRepo: transactionRepo,
	}
}

func (p *OutboxProcess) Start() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Println("ticker tickering")
		err := p.WriteTransactions()
		if err != nil {
			fmt.Println("outbox processor error:", err)
		}
	}
}

// update 10 pending transactions per call
// each update involves
// - write transaction into dynamodb
// - update outbox status in postgres
func (p *OutboxProcess) WriteTransactions() error {
	ctx := context.Background()
	pendingTransactions, err := p.outboxRepo.FetchPending(ctx, 10)
	if err != nil {
		return fmt.Errorf("Unable to do fetch transactions %w", err)
	}

	for _, transaction := range pendingTransactions {
		// update transactions
		_, err := p.transactionRepo.Create(ctx, &models.Transaction{
			TransactionID: transaction.ID,
			Side:          transaction.Side,
			Amount:        transaction.Amount,
			Description:   transaction.Description,
			CounterPartID: transaction.CounterPartID,
			AccountID:     transaction.AccountID,
			CreatedAt:     transaction.CreatedAt,
		})
		if err != nil {
			// if this fails, we can rerun this without problem
			fmt.Println("unable to write dynamodb transaction")
		} else {
			transaction.Status = models.OutboxTransactionStatusSETTLED
			transaction.RetryCount += 1
			err = p.outboxRepo.Update(ctx, *transaction)

			if err != nil {
				// if this fails, we can rerun this without problem
				fmt.Println("unable to update outbox transaction")
			}
		}
	}

	return nil
}
