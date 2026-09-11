-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS products (
  id UUID PRIMARY KEY,  -- using custom uuid type for better sqlc mapping
  name TEXT NOT NULL,
  price INTEGER NOT NULL
);

-- Add test Product
INSERT INTO products (id, name, price) VALUES ('00000000-0000-0000-0000-000000000001', 'kola', 100);
INSERT INTO products (id, name, price) VALUES ('00000000-0000-0000-0000-000000000002', 'fanta', 100);
INSERT INTO products (id, name, price) VALUES ('00000000-0000-0000-0000-000000000003', 'wasser', 70);
INSERT INTO products (id, name, price) VALUES ('00000000-0000-0000-0000-000000000004', 'sprite', 100);
INSERT INTO products (id, name, price) VALUES ('00000000-0000-0000-0000-000000000005', 'mezzo-mix', 100);
INSERT INTO products (id, name, price) VALUES ('00000000-0000-0000-0000-000000000006', 'spezi', 100);

-- +goose StatementEnd