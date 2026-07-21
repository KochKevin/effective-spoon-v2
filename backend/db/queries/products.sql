-- name: GetAllProducts :many
SELECT * FROM products;

-- name: GetProduct :one
SELECT 
id,
name,
price,
code
FROM products WHERE id = ?
ORDER BY rowid ASC;


-- name: GetProductByCode :one
SELECT 
id,
name,
price,
code
FROM products WHERE code = ?
ORDER BY rowid ASC;