-- +goose Up
-- +goose StatementBegin

-- LastStripeEvent
CREATE TABLE IF NOT EXISTS stripe_last_event_id(
    id INT PRIMARY KEY,
    event_id TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS chargements(
    id UUID PRIMARY KEY,
    status TEXT NOT NULL,
    user_id UUID NOT NULL,
    amount int NOT NULL,
    stripe_checkout_id TEXT,
    transaction_id UUID NOT NULL,

    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(transaction_id) REFERENCES transactions(id)
);

-- +goose StatementEnd