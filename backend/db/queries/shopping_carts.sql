-- name: CreateShoppingCart :one
INSERT INTO shopping_carts (id) VALUES (?) RETURNING *;

-- name: GetShoppingCart :one
SELECT
id
FROM shopping_carts
WHERE id = ?;

-- name: GetLineItemsOfShoppingCart :many
SELECT 
rel_shopping_carts_products.shopping_cart_id AS 'shoppingCartId', 
products.id AS 'productId',
products.name AS 'productName',
products.price AS 'productPrice', 
rel_shopping_carts_products.amount AS 'amount'

FROM rel_shopping_carts_products JOIN 
products ON products.id = rel_shopping_carts_products.product_id 
WHERE rel_shopping_carts_products.shopping_cart_id = ?;


-- name: DeleteAllLineItemsOfShoppingCart :exec
DELETE FROM rel_shopping_carts_products WHERE rel_shopping_carts_products.shopping_cart_id = ?;

-- name: CreateShoppingCartLineItem :exec
INSERT INTO rel_shopping_carts_products (
    shopping_cart_id,
    product_id,
    amount
) VALUES (?,?,?)