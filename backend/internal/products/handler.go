package products

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/KochKevin/effective-spoon-v2/internal/infrastructure"
	productsapi "github.com/KochKevin/effective-spoon-v2/internal/products/generated"
	"github.com/go-chi/render"
	"log/slog"
)

type Repo interface {
	GetProducts(ctx context.Context, tx *sql.Tx) ([]Product, error)
}

type Api struct {
	Repo Repo
	Txm  infrastructure.TxManager
}

// Load all products
// (GET /products)
func (a *Api) GetProducts(w http.ResponseWriter, r *http.Request) {

	//log.Println("Hallo projjjructs")

	var dtos []productsapi.Product

	err := a.Txm.WithTx(context.Background(), func(tx *sql.Tx) error {
		products, err := a.Repo.GetProducts(r.Context(), tx)

		//log.Println(products)

		if err != nil {
			//TODO: give client more inforamtion instead of an timeout
			slog.Error("Error loading products", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
		slog.Error("Error in /products transaction", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//log.Println(dtos)

	render.JSON(w, r, dtos)
}
