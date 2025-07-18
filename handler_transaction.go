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

// handlerMakeTransaction handles a POST api/send endpoint.
// Takes json with fields "from", "to", "amount" as a request body
// and returns json representation of created transaction in response
func (apiCFG *apiConfig) handlerMakeTransaction(w http.ResponseWriter, r *http.Request) {

	// parameters of JSON request body
	type parameters struct {
		From   uuid.UUID `json:"from"`
		To     uuid.UUID `json:"to"`
		Amount string    `json:"amount"`
	}

	// decoding JSON request body into params
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	// check if sender address == recipient address
	if params.From == params.To {
		respondWithError(w, http.StatusBadRequest, "Sender and recipient addresses can`t be the same")
		return
	}

	// check if amount < minimalTransactionAmount
	amount, err := decimal.NewFromString(params.Amount)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing amount: %v", err))
		return
	}
	if amount.LessThan(apiCFG.minimalTransactionAmount) {
		respondWithError(w, http.StatusBadRequest,
			fmt.Sprintf("Amout value is too small: minimum value is %v, your amount %v", apiCFG.minimalTransactionAmount, amount),
		)
		return
	}

	// check if sender balance < amount
	balanceStr, err := apiCFG.DB.GetWalletBalance(r.Context(), params.From)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn`t find wallet: %v", err))
		return
	}
	balance, err := decimal.NewFromString(balanceStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing balance: %v", err))
		return
	}
	if amount.GreaterThan(balance) {
		respondWithError(w, http.StatusConflict, fmt.Sprintf("Not enough money: balance is %v, your amount %v", balance, amount))
		return
	}

	// getting recipient balance
	recipientBalanceStr, err := apiCFG.DB.GetWalletBalance(r.Context(), params.To)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn`t find wallet: %v", err))
		return
	}

	recipientBalance, err := decimal.NewFromString(recipientBalanceStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing balance: %v", err))
		return
	}

	// new balance computations
	senderNewBalance := balance.Sub(amount)
	recipientNewBalance := recipientBalance.Add(amount)

	// making db transactions as atomic operation
	tx, err := apiCFG.dbConn.BeginTx(r.Context(), nil)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to begin transaction")
		return
	}

	defer tx.Rollback()

	// changing sender balance
	_, err = apiCFG.DB.WithTx(tx).ChangeWalletBalance(r.Context(), database.ChangeWalletBalanceParams{
		Balance:   senderNewBalance.String(),
		UpdatedAt: time.Now().UTC(),
		Address:   params.From,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error changing sender balance: %v", err))
		return
	}

	// changing recipient balance
	_, err = apiCFG.DB.WithTx(tx).ChangeWalletBalance(r.Context(), database.ChangeWalletBalanceParams{
		Balance:   recipientNewBalance.String(),
		UpdatedAt: time.Now().UTC(),
		Address:   params.To,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error changing receiver balance: %v", err))
		return
	}

	// creating transaction in db`s transactions table
	transaction, err := apiCFG.DB.WithTx(tx).AddTransaction(r.Context(), database.AddTransactionParams{
		ID:               uuid.New(),
		ExecutedAt:       time.Now().UTC(),
		Amount:           params.Amount,
		SenderAddress:    params.From,
		RecipientAddress: params.To,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating transaction: %v", err))
		return
	}

	// Commiting transaction
	if err := tx.Commit(); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	respondWithJSON(w, http.StatusCreated, dbTransactionToTransaction(transaction))

}

// handlerGetNLastTransactions handles a GET  /api/transactions?count=N endpoint.
// Returns json array of last N transactions as a response
func (apiCFG *apiConfig) handlerGetNLastTransactions(w http.ResponseWriter, r *http.Request) {
	// parsing count from query string
	countStr := r.URL.Query().Get("count")
	count, err := strconv.Atoi(countStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error converting count: %v", err))
		return
	}
	transactions, err := apiCFG.DB.GetNLastTransactions(r.Context(), int32(count))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting transactions: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, dbTransactionsToTransactions(transactions))
}
