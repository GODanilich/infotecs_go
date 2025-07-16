package main

import (
	"database/sql"
	"fmt"
	"infotecs_go/internal/database"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"

	_ "github.com/lib/pq"
)

// API config
type apiConfig struct {
	DB                       *database.Queries
	dbConn                   *sql.DB
	minimalTransactionAmount decimal.Decimal
}

func main() {
	// load environment variables from .env
	godotenv.Load(".env")

	// getting MINIMAL_TRANSACTION_AMOUNT from .env
	minimalTransactionAmountString := os.Getenv("MINIMAL_TRANSACTION_AMOUNT")
	if minimalTransactionAmountString == "" {
		log.Fatal("MINIMAL_TRANSACTION_AMOUNT is not found in the environment")
	}

	// Converting minimal transaction amount to decimal.Decimal
	minimumTransactionAmount, err := decimal.NewFromString(minimalTransactionAmountString)
	if err != nil {
		panic(err)
	}

	log.Printf("MINIMAL_TRANSACTION_AMOUNT is %v\n\n", minimumTransactionAmount)

	// getting PORT from .env
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	// getting DB_URL from .env
	db_URL := os.Getenv("DB_URL")
	if db_URL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	// connecting to db
	conn, err := sql.Open("postgres", db_URL)
	if err != nil {
		log.Fatal("Can`t connect to database:", err)
	}

	defer conn.Close()

	db := database.New(conn)

	apiCFG := apiConfig{
		DB:                       db,
		dbConn:                   conn,
		minimalTransactionAmount: minimumTransactionAmount,
	}

	// routing conf
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/api/wallet/{address}/balance", apiCFG.handlerGetBalance)
	v1Router.Post("/api/send", apiCFG.handlerMakeTransaction)
	v1Router.Get("/api/transactions", apiCFG.handlerGetNLastTransactions)

	router.Mount("/v1", v1Router)

	// check if wallets exists
	checkWallets(apiCFG)

	// configuring HTTP server
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Printf("Server is starting on port %v", portString)
	// starting the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("PORT is:", portString)
}
