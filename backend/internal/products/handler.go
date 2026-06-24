package products

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/KochKevin/effective-spoon-v2/internal/infrastructure"
	productsapi "github.com/KochKevin/effective-spoon-v2/internal/products/generated"
	"github.com/go-chi/render"
)

type Repo interface {
	GetProducts(ctx context.Context, tx *sql.Tx) ([]Product, error)
}

type Api struct {
	repo Repo
	txm  infrastructure.TxManager
}

// Load all products
// (GET /products)
func (a *Api) GetProducts(w http.ResponseWriter, r *http.Request) {

	var dtos []productsapi.Product

	err := a.txm.WithTx(context.Background(), func(tx *sql.Tx) error {
		products, err := a.repo.GetProducts(r.Context(), tx)

		if err != nil {
			//TODO: give client more inforamtion instead of an timeout
			return err
		}

		for _, p := range products {
			dtos = append(dtos, productsapi.Product{
				Id:    p.Id.String(),
				Name:  p.Name,
				Price: p.Price,
			})
		}

		return nil

	})

	if err != nil {
		return
	}

	render.JSON(w, r, dtos)
}
