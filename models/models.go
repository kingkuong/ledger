package models

type Account struct {
	accountID     string
	accountNumber string
	nickname      *string
	balance       int
	createdAt     string
	updatedAt     string
}

type Transaction struct {
	transactionID string
	side          string
	amount        int
	description   string
	counterPartID string
	accountID     string
	createdAt     string
	updatedAt     string
}

type OutboxTransaction struct {
	*Transaction
	status      string
	retryCount  int
	lastRetryAt string
	updatedAt   string
	createdAt   string
}
