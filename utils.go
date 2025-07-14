package main

import (
	"context"
	"infotecs_go/internal/database"
	"time"

	"github.com/google/uuid"
)

// checkWallets generates 10 wallets
// with 100 on balance if db is empty
func checkWallets(apiCFG apiConfig) {
	if wallets, _ := apiCFG.DB.GetWallets(context.Background()); len(wallets) == 0 {
		for i := 0; i < 10; i++ {
			apiCFG.DB.CreateWallet(context.Background(), database.CreateWalletParams{
				Address:   uuid.New(),
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
				Balance:   "100",
			})
		}
	}
}
