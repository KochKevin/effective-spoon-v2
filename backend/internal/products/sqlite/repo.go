package productssqlite

import (
	"context"
	"database/sql"

	sqlc "github.com/KochKevin/effective-spoon-v2/internal/infrastructure/sqlite/generated"
	"github.com/KochKevin/effective-spoon-v2/internal/products"
)

type Repo struct {
	Queries sqlc.Queries
}

//func NewRepo(db)

func (r *Repo) GetProducts(ctx context.Context, tx *sql.Tx) (domains []products.Product, err error) {

	dbProducts, err := r.Queries.WithTx(tx).GetAllProducts(ctx)
	if err != nil {
		return nil, err
	}

	for _, product := range dbProducts {

		domains = append(domains, products.NewProduct(product.ID, product.Name, int(product.Price)))
	}

	return domains, nil
}
