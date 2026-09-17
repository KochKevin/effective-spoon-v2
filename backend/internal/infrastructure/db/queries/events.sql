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

-- name: UpsertEventUsage :one
INSERT INTO event_usage(
    user_id,
    event_id,
    used_amount_free_products
)
VALUES(?,?,?)
ON CONFLICT(user_id, event_id) 
DO 
   UPDATE SET used_amount_free_products = excluded.used_amount_free_products

RETURNING *;

-- name: GetEventUsage :one
SELECT 
user_id,
event_id,
used_amount_free_products
FROM event_usage 
WHERE user_id = ? AND event_id = ?;
