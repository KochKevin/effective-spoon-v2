-- name: CreateShoppingCart :one
INSERT INTO shopping_carts (id, user_id, transaction_id, status) VALUES (?, ?, ?, ?) RETURNING *;

-- name: GetShoppingCart :one
SELECT
id,
user_id,
transaction_id,
status
FROM shopping_carts
WHERE id = ?;

-- name: UpdateShoppingCart :exec
UPDATE shopping_carts 
SET 
user_id = ?,
transaction_id = ?,
status = ?
WHERE id = ?;



-- Line Items


-- name: GetLineItemsOfShoppingCart :many
SELECT 
rel_shopping_carts_products.shopping_cart_id AS 'shoppingCartId', 
products.id AS 'productId',
products.name AS 'productName',
products.price AS 'productPrice', 
rel_shopping_carts_products.amount AS 'amount'
FROM rel_shopping_carts_products JOIN 
products ON products.id = rel_shopping_carts_products.product_id 
WHERE rel_shopping_carts_products.shopping_cart_id = ?
ORDER BY rel_shopping_carts_products.rowid ASC; --order by to keep the order in which they are inserted


-- name: DeleteAllLineItemsOfShoppingCart :exec
DELETE FROM rel_shopping_carts_products WHERE rel_shopping_carts_products.shopping_cart_id = ?;

-- name: CreateShoppingCartLineItem :exec
INSERT INTO rel_shopping_carts_products (
    shopping_cart_id,
    product_id,
    amount
) VALUES (?,?,?)