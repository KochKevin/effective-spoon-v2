package shoppingcartssqlite

import (
	"context"
	"database/sql"

	"log/slog"

	sqlc "github.com/KochKevin/effective-spoon-v2/internal/infrastructure/sqlite/generated"
	"github.com/KochKevin/effective-spoon-v2/internal/shoppingcarts"
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

func (r *Repo) GetShoppingCart(ctx context.Context, tx *sql.Tx) (shoppingcarts.ShoppingCart, error) {
	panic("not implemented") // TODO: Implement
}

func (r *Repo) SaveShoppingCart(ctx context.Context, tx *sql.Tx, cart shoppingcarts.ShoppingCart) (shoppingcarts.ShoppingCart, error) {
	panic("not implemented") // TODO: Implement
}
