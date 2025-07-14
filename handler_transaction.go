package main

import (
	"encoding/json"
	"fmt"
	"infotecs_go/internal/database"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func (apiCFG *apiConfig) handlerMakeTransaction(w http.ResponseWriter, r *http.Request) {

	minimumTransactionAmount, err := decimal.NewFromString("0.01")
	if err != nil {
		panic(err)
	}
	type parameters struct {
		From   uuid.UUID `json:"from"`
		To     uuid.UUID `json:"to"`
		Amount string    `json:"amount"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	balanceStr, err := apiCFG.DB.GetWalletBalance(r.Context(), params.From)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn`t find wallet: %v", err))
		return
	}

	amount, err := decimal.NewFromString(params.Amount)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing amount: %v", err))
		return
	}
	balance, err := decimal.NewFromString(balanceStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing balance: %v", err))
		return
	}

	if amount.LessThan(minimumTransactionAmount) {
		respondWithError(w, 400, fmt.Sprintf("Amout value is too small: minimum value is %v, your amount %v", minimumTransactionAmount, amount))
		return
	}

	if amount.GreaterThan(balance) {
		respondWithError(w, 400, fmt.Sprintf("Not enough money: balance is %v, your amount %v", balance, amount))
		return
	}

	recieverBalanceStr, err := apiCFG.DB.GetWalletBalance(r.Context(), params.To)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn`t find wallet: %v", err))
		return
	}

	recieverBalance, err := decimal.NewFromString(recieverBalanceStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing balance: %v", err))
		return
	}

	senderNewBalance := balance.Sub(amount)
	recieverNewBalance := recieverBalance.Add(amount)
	_, err = apiCFG.DB.ChangeWalletBalance(r.Context(), database.ChangeWalletBalanceParams{
		Balance: senderNewBalance.String(),
		Address: params.From,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error changing sender balance: %v", err))
		return
	}
	_, err = apiCFG.DB.ChangeWalletBalance(r.Context(), database.ChangeWalletBalanceParams{
		Balance: recieverNewBalance.String(),
		Address: params.To,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error changing reciever balance: %v", err))
		return
	}

	transaction, err := apiCFG.DB.AddTransaction(r.Context(), database.AddTransactionParams{
		ID:               uuid.New(),
		ExecutedAt:       time.Now().UTC(),
		Amount:           params.Amount,
		SenderAddress:    params.From,
		RecipientAddress: params.To,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating transaction: %v", err))
		return
	}

	respondWithJSON(w, 200, transaction)
}

func (apiCFG *apiConfig) handlerGetNLastTransactions(w http.ResponseWriter, r *http.Request) {
	countStr := r.URL.Query().Get("count")
	count, err := strconv.Atoi(countStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error converting count: %v", err))
		return
	}
	transactions, err := apiCFG.DB.GetNLastTransactions(r.Context(), int32(count))
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error getting transactions: %v", err))
		return
	}

	respondWithJSON(w, 200, transactions)
}
