package shoppingcartssqlite

import (
	"context"
	"database/sql"

	"log/slog"

	sqlc "github.com/KochKevin/effective-spoon-v2/internal/infrastructure/sqlite/generated"
	"github.com/KochKevin/effective-spoon-v2/internal/money"
	"github.com/KochKevin/effective-spoon-v2/internal/products"
	"github.com/KochKevin/effective-spoon-v2/internal/shoppingcarts"
	"github.com/google/uuid"
)

type Repo struct {
	Queries sqlc.Queries
}

// TODO: calling Create should also create the line items
func (r *Repo) CreateShoppingCart(ctx context.Context, tx *sql.Tx, cart shoppingcarts.ShoppingCart) (shoppingcarts.ShoppingCart, error) {

	obj, err := r.Queries.WithTx(tx).CreateShoppingCart(ctx, sqlc.CreateShoppingCartParams{
		ID:            cart.Id,
		UserID:        cart.UserId,
		TransactionID: cart.TransactionId,
		Status:        string(cart.Status),
	})
	if err != nil {
		slog.Error("Error in products query", err)
		return shoppingcarts.ShoppingCart{}, err
	}

	return shoppingcarts.ShoppingCartFrom(obj.ID, nil, obj.UserID, obj.TransactionID, shoppingcarts.ShoppingCartStatus(obj.Status)), nil
}

func (r *Repo) GetShoppingCart(ctx context.Context, tx *sql.Tx, id uuid.UUID) (shoppingcarts.ShoppingCart, error) {
	shoppingCart, err := r.Queries.WithTx(tx).GetShoppingCart(ctx, id)
	if err != nil {
		slog.Error("Error in get shopping cart query", err)
		return shoppingcarts.ShoppingCart{}, err
	}

	lineItems, err := r.Queries.WithTx(tx).GetLineItemsOfShoppingCart(ctx, id)
	if err != nil {
		slog.Error("Error in get line items of shopping cart query", err)
		return shoppingcarts.ShoppingCart{}, err
	}

	var cart shoppingcarts.ShoppingCart
	cart.Id = shoppingCart.ID
	cart.UserId = shoppingCart.UserID
	cart.TransactionId = shoppingCart.TransactionID
	cart.Status = shoppingcarts.ShoppingCartStatus(shoppingCart.Status)

	for _, item := range lineItems {
		cart.LineItems = append(cart.LineItems, shoppingcarts.LineItem{
			Amount: int(item.Amount),
			Product: products.Product{
				Id:    item.Productid,
				Name:  item.Productname,
				Price: money.MoneyFrom(int(item.Productprice)),
			},
		})
	}

	return cart, nil
}

func (r *Repo) SaveShoppingCart(ctx context.Context, tx *sql.Tx, cart shoppingcarts.ShoppingCart) error {

	//Update Line Items
	err := r.Queries.WithTx(tx).DeleteAllLineItemsOfShoppingCart(ctx, cart.Id)
	if err != nil {
		slog.Error("Error in deleting all line items of an shopping cart query", err)
		return err
	}

	for _, item := range cart.LineItems {

		err := r.Queries.WithTx(tx).CreateShoppingCartLineItem(ctx, sqlc.CreateShoppingCartLineItemParams{
			ShoppingCartID: cart.Id,
			ProductID:      item.Product.Id,
			Amount:         int64(item.Amount),
		})
		if err != nil {
			slog.Error("Error in creating line item query", err)
			return err
		}
	}

	//Update Shopping Cart
	r.Queries.WithTx(tx).UpdateShoppingCart(ctx, sqlc.UpdateShoppingCartParams{
		UserID:        cart.UserId,
		TransactionID: cart.TransactionId,
		ID:            cart.Id,
		Status:        string(cart.Status),
	})

	return nil
}
