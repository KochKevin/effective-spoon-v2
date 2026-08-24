-- name: GetLastStripeEventId :one
SELECT event_id 
FROM stripe_last_event_id WHERE id = 1;

-- name: SetLastStripeEventId :exec
UPDATE stripe_last_event_id 
SET event_id = ?
WHERE id = 1;


-- name: GetChargementIntent :one
SELECT 
id, 
status, 
user_id, 
amount, 
stripe_checkout_id, 
transaction_id 
FROM chargements WHERE id = ?;

-- name: UpdateChargementIntent :exec
UPDATE chargements
SET
status = ?,
user_id = ?,
amount = ?, 
stripe_checkout_id = ?,
transaction_id = ? 
WHERE id = ?;


-- name: CreateChargementIntent :one
INSERT INTO 
chargements 
(id, status, user_id, amount, stripe_checkout_id, transaction_id) 
VALUES (?, ?, ?, ?, ?, ?) RETURNING *;
