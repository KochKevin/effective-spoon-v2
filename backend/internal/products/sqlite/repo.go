package productssqlite

import (
	"context"

	"github.com/KochKevin/effective-spoon-v2/internal/products"
)

type Repo struct {
}

//func NewRepo(db)

func (r *Repo) GetProducts(ctx context.Context) []products.Product {
	panic("not implemented") // TODO: Implement
}
