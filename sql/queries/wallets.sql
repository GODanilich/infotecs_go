-- name: GetWallet :one
SELECT * FROM wallets WHERE address = $1;

-- name: GetWallets :many
SELECT * FROM wallets;

-- name: CreateWallet :one
INSERT INTO wallets (address, created_at, updated_at, balance)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetWalletBalance :one
SELECT balance FROM wallets WHERE address = $1;

-- name: ChangeWalletBalance :one
UPDATE wallets
SET
    balance = $1,
    updated_at = NOW()
WHERE address = $2
RETURNING *;