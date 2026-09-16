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

CREATE TABLE IF NOT EXISTS event_usage(
    user_id UUID NOT NULL,
    event_id UUID NOT NULL,
    used_amount_free_products INT NOT NULL,

    PRIMARY KEY(user_id, event_id),
    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(event_id) REFERENCES events(id)
);
-- +goose StatementEnd