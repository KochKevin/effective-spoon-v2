-- +goose Up
-- +goose StatementBegin


CREATE TABLE IF NOT EXISTS events(
    id UUID PRIMARY KEY,
    -- user_id is the id of the event creator
    user_id UUID NOT NULL,
    amount_free_products_per_user INT NOT NULL,
    start_timestamp DATETIME NOT NULL,
    end_timestamp DATETIME NOT NULL,
    status TEXT NOT NULL,

    FOREIGN KEY(user_id) REFERENCES users(id)
);

-- ensure that only one event can be active at once
CREATE UNIQUE INDEX idx_unique_active_event 
ON events(status) 
WHERE status = "active";


ALTER TABLE shopping_carts ADD COLUMN use_event BOOL NOT NULL;
ALTER TABLE shopping_carts ADD COLUMN event_id UUID NOT NULL REFERENCES events(id);

ALTER TABLE rel_shopping_carts_products ADD COLUMN amount_free_products INT NOT NULL;
-- +goose StatementEnd