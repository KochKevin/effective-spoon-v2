package products

import (
	"context"
	"net/http"

	productsapi "github.com/KochKevin/effective-spoon-v2/internal/products/generated"
	"github.com/go-chi/render"
)

type ProductRepository interface {
	GetProducts(ctx context.Context) []Product
}

type Api struct {
	repo ProductRepository
}

// Load all products
// (GET /products)
func (a *Api) GetProducts(w http.ResponseWriter, r *http.Request) {

	products := a.repo.GetProducts(r.Context())

	var dtos []productsapi.Product

	for _, p := range products {
		dtos = append(dtos, productsapi.Product{
			Id:    p.Id.String(),
			Name:  p.Name,
			Price: p.Price,
		})
	}

	render.JSON(w, r, dtos)
}
