-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS products (
  id UUID PRIMARY KEY,  -- using custom uuid type for better sqlc mapping
  name TEXT NOT NULL,
  price INTEGER NOT NULL
);

-- +goose StatementEnd