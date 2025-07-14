-- +goose Up

CREATE TABLE wallets(
    address UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    balance NUMERIC(10,2) NOT NULL
);

-- +goose Down
DROP TABLE wallets;