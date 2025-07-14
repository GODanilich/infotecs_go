-- name: AddTransaction :one
INSERT INTO transactions (id, executed_at, amount, sender_address, recipient_address)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetNLastTransactions :many
SELECT * FROM transactions
ORDER BY executed_at DESC
LIMIT $1;