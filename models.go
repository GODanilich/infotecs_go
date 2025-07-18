package main

import (
	"infotecs_go/internal/database"
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID               uuid.UUID `json:"id"`
	ExecutedAt       time.Time `json:"executed_at"`
	Amount           string    `json:"amount"`
	SenderAddress    uuid.UUID `json:"sender_address"`
	RecipientAddress uuid.UUID `json:"recipient_address"`
}

type Wallet struct {
	Address   uuid.UUID `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Balance   string    `json:"balance"`
}

func dbTransactionToTransaction(dbTransaction database.Transaction) Transaction {
	return Transaction{
		ID:               dbTransaction.ID,
		ExecutedAt:       dbTransaction.ExecutedAt,
		Amount:           dbTransaction.Amount,
		SenderAddress:    dbTransaction.SenderAddress,
		RecipientAddress: dbTransaction.RecipientAddress,
	}
}

func dbTransactionsToTransactions(dbTransactions []database.Transaction) (transactions []Transaction) {
	for _, dbTransaction := range dbTransactions {
		transactions = append(transactions, dbTransactionToTransaction(dbTransaction))
	}
	return transactions
}
