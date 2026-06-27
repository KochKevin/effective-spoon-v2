-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS products (
  id UUID PRIMARY KEY,  -- using custom uuid type for better sqlc mapping
  name TEXT NOT NULL,
  price INTEGER NOT NULL
);

-- Add test Product
INSERT INTO products (id, name, price) VALUES ('00000000-0000-0000-0000-000000000000', 'test', 100)

-- +goose StatementEnd