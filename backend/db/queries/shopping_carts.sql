-- name: CreateShoppingCart :one
INSERT INTO shopping_carts (id) VALUES (?) RETURNING *;