package shoppingcarts

import (
	"context"
	"database/sql"
	"net/http"

	"log/slog"

	"github.com/KochKevin/effective-spoon-v2/internal/infrastructure"
	"github.com/KochKevin/effective-spoon-v2/internal/products"
	shoppingcartsapi "github.com/KochKevin/effective-spoon-v2/internal/shoppingcarts/generated"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type Repo interface {
	CreateShoppingCart(ctx context.Context, tx *sql.Tx, cart ShoppingCart) (ShoppingCart, error)
	GetShoppingCart(ctx context.Context, tx *sql.Tx, id uuid.UUID) (ShoppingCart, error)
	SaveShoppingCart(ctx context.Context, tx *sql.Tx, cart ShoppingCart) error
}

type ProductRepo interface {
	GetProduct(ctx context.Context, tx *sql.Tx, id uuid.UUID) (products.Product, error)
}

type Api struct {
	Repo Repo
	ProductRepo ProductRepo
	Txm  infrastructure.TxManager
}

func (a *Api) ToDto(cart ShoppingCart) shoppingcartsapi.ShoppingCart {

	var lineItems []shoppingcartsapi.LineItem

	for _, item := range cart.LineItems {

		lineItems = append(lineItems, shoppingcartsapi.LineItem{
			Amount:    item.Amount,
			Price:     float32(item.GetPrice().GetAsEuro()),
			ProductId: item.Product.Id.String(),
			ProductName: item.Product.Name,
		})
	}

	return shoppingcartsapi.ShoppingCart{
		Id:        cart.Id.String(),
		FullPrice: float32(cart.GetFullPrice().GetAsEuro()),
		LineItems: lineItems,
	}

}

// Create a new shopping cart
// (POST /shopping-carts)
func (a *Api) PostShoppingCarts(w http.ResponseWriter, r *http.Request) {

	var dto shoppingcartsapi.ShoppingCart

	err := a.Txm.WithTx(context.Background(), func(tx *sql.Tx) error {

		cart, err := a.Repo.CreateShoppingCart(r.Context(), tx, NewShoppingCart())

		if err != nil {
			//TODO: give client more inforamtion instead of an timeout
			slog.Error("Error creating shopping cart", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return err
		}

		dto = a.ToDto(cart)

		return nil

	})

	if err != nil {
		slog.Error("Error in /products transaction", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, dto)
}

// Remove product from shopping cart
// (POST /shopping-carts/{id}/decrease)
func (a *Api) PostShoppingCartsIdDecrease(w http.ResponseWriter, r *http.Request, id string, params shoppingcartsapi.PostShoppingCartsIdDecreaseParams) {
	panic("not implemented") // TODO: Implement
}

// Add product to shopping cart
// (POST /shopping-carts/{id}/increase)
func (a *Api) PostShoppingCartsIdIncrease(w http.ResponseWriter, r *http.Request, id string, params shoppingcartsapi.PostShoppingCartsIdIncreaseParams) {
	slog.Debug("Increase Product by id")

	var dto shoppingcartsapi.ShoppingCart

	err := a.Txm.WithTx(context.Background(), func(tx *sql.Tx) error {

		cart, err := a.Repo.GetShoppingCart(r.Context(), tx, uuid.MustParse(id))

		slog.Debug("Before increase: ", cart)

		if err != nil {
			//TODO: give client more inforamtion instead of an timeout
			slog.Error("Error getting shopping cart", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return err
		}

		productToAdd, err := a.ProductRepo.GetProduct(r.Context(), tx, uuid.MustParse(params.ProductID))
		if err != nil {
			//TODO: give client more inforamtion instead of an timeout
			slog.Error("Error getting product to add to shopping cart", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return err
		}


		cart.IncreaseProductAmount(productToAdd)

		err = a.Repo.SaveShoppingCart(r.Context(), tx, cart)
		if err != nil {
			//TODO: give client more inforamtion instead of an timeout
			slog.Error("Error getting shopping cart", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return err
		}

		//slog.Debug("After increase: ", "cart", cart)

		dto = a.ToDto(cart)
		slog.Debug("After increase: ", "dto", dto)

		return nil

	})

	if err != nil {
		slog.Error("Error in /products transaction", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, dto)
}

// Load all products
// (GET /products)
func (a *Api) GetProducts(w http.ResponseWriter, r *http.Request) {

	//log.Println("Hallo projjjructs")
	/*
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
	*/
}
