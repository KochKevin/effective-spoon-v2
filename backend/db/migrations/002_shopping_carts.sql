-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS shopping_carts (
  id UUID PRIMARY KEY  -- using custom uuid type for better sqlc mapping
);

CREATE TABLE IF NOT EXISTS rel_shopping_carts_products (
  shopping_cart_id UUID NOT NULL,
  product_id UUID NOT NULL,
  amount INT NOT NULL,
  FOREIGN KEY(shopping_cart_id) REFERENCES shopping_carts(id),
  FOREIGN KEY(product_id) REFERENCES products(id),
  PRIMARY KEY(shopping_cart_id, product_id)
);

-- +goose StatementEnd