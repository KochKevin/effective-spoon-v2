-- +goose Up
-- +goose StatementBegin

-- Users
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY,  -- using custom uuid type for better sqlc mapping
  name TEXT NOT NULL
);

-- Add test user
INSERT INTO users (id, name) VALUES ('00000000-0000-0000-0000-000000000001', 'nico');
INSERT INTO users (id, name) VALUES ('00000000-0000-0000-0000-000000000002', 'max mustermann');
INSERT INTO users (id, name) VALUES ('00000000-0000-0000-0000-000000000003', 'max muster');

-- Update shopping_carts table, so that its owned by an user
ALTER TABLE shopping_carts ADD COLUMN user_id UUID NOT NULL REFERENCES users(id);




-- +goose StatementEnd