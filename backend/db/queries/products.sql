-- name: GetAllProducts :many
SELECT * FROM products;

-- name: GetProduct :one
SELECT 
id,
name,
price
FROM products WHERE id = ?;