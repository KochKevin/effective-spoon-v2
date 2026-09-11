-- +goose Up
-- +goose StatementBegin

-- Transactions
CREATE TABLE IF NOT EXISTS user_transactions(
    id UUID PRIMARY KEY,  -- using custom uuid type for better sqlc mapping
    user_id UUID NOT NULL,
    amount INT NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id)
);

-- Update shopping_carts table, so that it references to an transaction
ALTER TABLE shopping_carts ADD COLUMN transaction_id UUID REFERENCES user_transactions(id);
ALTER TABLE shopping_carts ADD COLUMN status text NOT NULL;

-- +goose StatementEnd