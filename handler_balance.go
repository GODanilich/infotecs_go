package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// handlerGetBalance handles GET /api/wallet/{address}/balance endpoint
func (apiCFG *apiConfig) handlerGetBalance(w http.ResponseWriter, r *http.Request) {
	addressStr := chi.URLParam(r, "address")
	address, err := uuid.Parse(addressStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn`t parse wallet`s address: %v", err))
		return
	}
	balance, err := apiCFG.DB.GetWalletBalance(r.Context(), address)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Couldn`t find wallet: %v", err))
		return
	}
	respondWithJSON(w, http.StatusOK, balance)
}
