package main

import (
	"context"
	"fake-ledger/models"
	service "fake-ledger/services"
	"fmt"
)

func main() {
	fmt.Println("hello")
	ctx := context.Background()
	s, err := service.NewLedgerService(ctx)
	if err != nil {
		panic(fmt.Errorf("Unable to start the service, check your code!! %w", err))
	}

	fmt.Println("service registered")
	account := &models.Account{
		AccountNumber: "ACC0001",
	}
	acct, err := s.CreateAccount(ctx, account)
	if err != nil {
		panic(fmt.Errorf("Unable to create account!! %w", err))
	}

	fmt.Printf("----------ACCOUNT----------\n%+v\n", acct)

	transaction := &models.Transaction{
		Side:        "DR",
		Amount:      5000,
		Description: "purchase",
		AccountID:   acct.ID,
	}
	trans, err := s.CreateTransaction(ctx, transaction)
	if err != nil {
		panic(fmt.Errorf("Unable to create transaction%w", err))
	}

	fmt.Printf("----------TRANSACTION----------\n%+v\n", trans)
}
