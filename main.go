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
	account, err = s.CreateAccount(ctx, account)
	if err != nil {
		panic(fmt.Errorf("Unable to create account!! %w", err))
	}
	accountTo := &models.Account{
		AccountNumber: "ACC0002",
	}
	accountTo, err = s.CreateAccount(ctx, accountTo)
	if err != nil {
		panic(fmt.Errorf("Unable to create account!! %w", err))
	}

	fmt.Printf("----------ACCOUNTs----------\n%+v\n%+v\n", account, accountTo)

	transaction := &models.Transaction{
		Side:        "CR",
		Amount:      5000,
		Description: "Initial Credit",
		AccountID:   account.ID,
	}
	transaction, err = s.CreateTransaction(ctx, transaction)
	if err != nil {
		panic(fmt.Errorf("Unable to create transaction%w", err))
	}

	fmt.Printf("----------TRANSACTION----------\n%+v\n", transaction)

	err = s.Transfer(ctx, account.ID, accountTo.ID, 100)
	if err != nil {
		panic(fmt.Errorf("Unable to do transfer%w", err))
	}

	fmt.Printf("----------TRANSFER----------\n%+v\n", transaction)
}
