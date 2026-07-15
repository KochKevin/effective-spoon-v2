-- name: GetUser :one
SELECT 
users.id,
users.name,
CAST(COALESCE(
(
    SELECT SUM(user_transactions.amount) 
    FROM user_transactions 
    WHERE user_transactions.user_id = users.id
),0
) AS INTEGER) as balance -- Fallback to zero if no transaction exists
FROM users
WHERE users.id = ?;


-- name: CreateTransaction :one
INSERT INTO user_transactions (
id,
user_id,
amount
)
VALUES (?,?,?) RETURNING *;