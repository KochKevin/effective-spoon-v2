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

func (r *Repo) CreateShoppingCart(ctx context.Context, tx *sql.Tx, cart shoppingcarts.ShoppingCart) (shoppingcarts.ShoppingCart, error) {

	obj, err := r.Queries.WithTx(tx).CreateShoppingCart(ctx, cart.Id)
	if err != nil {
		slog.Error("Error in products query", err)
		return shoppingcarts.ShoppingCart{}, err
	}

	return shoppingcarts.ShoppingCartFrom(obj, nil), nil
}

func (r *Repo) GetShoppingCart(ctx context.Context, tx *sql.Tx, id uuid.UUID) (shoppingcarts.ShoppingCart, error) {
	data, err := r.Queries.WithTx(tx).GetShoppingCart(ctx, id)
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

	cart.Id = data

	return cart, nil
}

func (r *Repo) SaveShoppingCart(ctx context.Context, tx *sql.Tx, cart shoppingcarts.ShoppingCart) error {

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

	return nil
}
