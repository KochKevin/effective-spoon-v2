-- +goose Up
-- +goose StatementBegin

-- Transactions
CREATE TABLE IF NOT EXISTS user_transactions(
    id UUID PRIMARY KEY,  -- using custom uuid type for better sqlc mapping
    user_id UUID NOT NULL,
    amount INT NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id)
);



-- +goose StatementEnd