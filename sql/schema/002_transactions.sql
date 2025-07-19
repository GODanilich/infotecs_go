-- +goose Up

CREATE TABLE transactions(
    id UUID PRIMARY KEY,
    executed_at TIMESTAMP NOT NULL,
    amount NUMERIC(15,2) NOT NULL,
    sender_address UUID NOT NULL REFERENCES wallets(address) ON DELETE CASCADE,
    recipient_address UUID NOT NULL REFERENCES wallets(address) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE transactions;