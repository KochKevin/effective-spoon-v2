-- name: CreateEvent :one
INSERT INTO events (
id,
user_id,
amount_free_products_per_user,
start_timestamp,
end_timestamp,
status
)
VALUES (?,?,?,?,?,?) RETURNING *;



-- name: GetEvent :one
SELECT
id,
user_id,
amount_free_products_per_user,
start_timestamp,
end_timestamp,
status
FROM events
WHERE id = ?;

-- name: GetEventByStatus :one
SELECT
id,
user_id,
amount_free_products_per_user,
start_timestamp,
end_timestamp,
status
FROM events
WHERE status = ?
LIMIT 1;


-- name: GetEventUsage :one
SELECT 
CAST( COALESCE(SUM(rel_shopping_carts_products.amount_free_products), 0) AS INTEGER) AS event_usage
FROM
rel_shopping_carts_products
JOIN shopping_carts ON shopping_carts.id = rel_shopping_carts_products.shopping_cart_id
WHERE shopping_carts.event_id = ? 
AND shopping_carts.user_id = ?
AND shopping_carts.use_event = ?;