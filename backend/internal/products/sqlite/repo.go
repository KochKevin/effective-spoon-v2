package productssqlite

import (
	"context"
	"database/sql"

	"log/slog"

	sqlc "github.com/KochKevin/effective-spoon-v2/internal/infrastructure/sqlite/generated"
	"github.com/KochKevin/effective-spoon-v2/internal/money"
	"github.com/KochKevin/effective-spoon-v2/internal/products"
	"github.com/google/uuid"
)

type Repo struct {
	Queries sqlc.Queries
}

//func NewRepo(db)

func (r *Repo) GetProducts(ctx context.Context, tx *sql.Tx) (domains []products.Product, err error) {

	dbProducts, err := r.Queries.WithTx(tx).GetAllProducts(ctx)
	if err != nil {
		slog.Error("Error in products query", err)
		return nil, err
	}

	for _, product := range dbProducts {

		domains = append(domains, products.NewProduct(product.ID, product.Name, money.MoneyFrom(int(product.Price))))
	}

	return domains, nil
}

func (r *Repo) GetProduct(ctx context.Context, tx *sql.Tx, id uuid.UUID) (products.Product, error) {

	product, err := r.Queries.WithTx(tx).GetProduct(ctx, id)
	if err != nil {
		slog.Error("Error in products query", err)
		return products.Product{}, err
	}

	return products.NewProduct(product.ID, product.Name, money.MoneyFrom(int(product.Price))), nil

}

func (r *Repo) GetProductByCode(ctx context.Context, tx *sql.Tx, code string) (products.Product, error) {

	product, err := r.Queries.WithTx(tx).GetProductByCode(ctx, code)
	if err != nil {
		slog.Error("Error in get products by code query", "error", err)
		return products.Product{}, err
	}

	return products.NewProduct(product.ID, product.Name, money.MoneyFrom(int(product.Price))), nil

}
