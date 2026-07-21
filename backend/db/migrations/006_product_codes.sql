-- +goose Up
-- +goose StatementBegin

-- NEEDS TO BE CHANGED LATER 
ALTER TABLE products ADD COLUMN code text NOT NULL DEFAULT '00000';

UPDATE products SET code = 1 WHERE products.id = '00000000-0000-0000-0000-000000000001';
UPDATE products SET code = 2 WHERE products.id = '00000000-0000-0000-0000-000000000002';
UPDATE products SET code = 3 WHERE products.id = '00000000-0000-0000-0000-000000000003';
UPDATE products SET code = 4 WHERE products.id = '00000000-0000-0000-0000-000000000004';
UPDATE products SET code = 5 WHERE products.id = '00000000-0000-0000-0000-000000000005';
UPDATE products SET code = 6 WHERE products.id = '00000000-0000-0000-0000-000000000006';

-- +goose StatementEnd