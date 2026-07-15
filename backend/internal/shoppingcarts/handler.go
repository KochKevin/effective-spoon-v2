package shoppingcarts

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"log/slog"

	"github.com/KochKevin/effective-spoon-v2/internal/infrastructure"
	"github.com/KochKevin/effective-spoon-v2/internal/products"
	shoppingcartsapi "github.com/KochKevin/effective-spoon-v2/internal/shoppingcarts/generated"
	"github.com/KochKevin/effective-spoon-v2/internal/users"
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

type UserRepo interface {
	CreateTransaction(ctx context.Context, tx *sql.Tx, transaction users.Transaction) (users.Transaction, error)
}

type Api struct {
	Repo        Repo
	ProductRepo ProductRepo
	UserRepo    UserRepo
	Txm         infrastructure.TxManager
}

//a *Api github.com/KochKevin/effective-spoon-v2/internal/shoppingcarts/generated.ServerInterface

func (a *Api) ToDto(cart ShoppingCart) shoppingcartsapi.ShoppingCart {

	var lineItems []shoppingcartsapi.LineItem

	for _, item := range cart.LineItems {

		lineItems = append(lineItems, shoppingcartsapi.LineItem{
			Amount:      item.Amount,
			Price:       float32(item.GetPrice().GetAsEuro()),
			ProductId:   item.Product.Id.String(),
			ProductName: item.Product.Name,
		})
	}

	return shoppingcartsapi.ShoppingCart{
		Id:        cart.Id.String(),
		FullPrice: float32(cart.GetFullPrice().GetAsEuro()),
		LineItems: lineItems,
		UserId:    cart.UserId.String(),
	}

}

// Create a new shopping cart
// (POST /shopping-carts)
func (a *Api) PostShoppingCarts(w http.ResponseWriter, r *http.Request) {

	var dto shoppingcartsapi.ShoppingCart

	err := a.Txm.WithTx(context.Background(), func(tx *sql.Tx) error {

		userID, ok := r.Context().Value("user_id").(uuid.UUID)
		if !ok {
			//slog.Error("Can not get user_id from context")
			return errors.New("Can not get user_id from context")
		}

		cart, err := a.Repo.CreateShoppingCart(r.Context(), tx, NewShoppingCart(userID))

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
	slog.Debug("Decrease Product by id")

	var dto shoppingcartsapi.ShoppingCart

	err := a.Txm.WithTx(context.Background(), func(tx *sql.Tx) error {

		cart, err := a.Repo.GetShoppingCart(r.Context(), tx, uuid.MustParse(id))

		slog.Debug("Before decrease: ", cart)

		if err != nil {
			//TODO: give client more inforamtion instead of an timeout
			slog.Error("Error getting shopping cart", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return err
		}

		cart.DecreaseProductAmount(uuid.MustParse(params.ProductID))

		err = a.Repo.SaveShoppingCart(r.Context(), tx, cart)
		if err != nil {
			//TODO: give client more inforamtion instead of an timeout
			slog.Error("Error getting shopping cart", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return err
		}

		dto = a.ToDto(cart)
		slog.Debug("After decrease: ", "dto", dto)

		return nil

	})

	if err != nil {
		slog.Error("Error in /shoppingcart decrease transaction", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, dto)

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

// Buy products in shopping cart
// (POST /shopping-carts/{id}/checkout)
func (a *Api) PostShoppingCartsIdCheckout(w http.ResponseWriter, r *http.Request, id string) {

	var dto shoppingcartsapi.ShoppingCart

	err := a.Txm.WithTx(context.Background(), func(tx *sql.Tx) error {

		/*
			userID, ok := r.Context().Value("user_id").(uuid.UUID)
			if !ok {
				//slog.Error("Can not get user_id from context")
				return errors.New("Can not get user_id from context")
			}
		*/

		cart, err := a.Repo.GetShoppingCart(r.Context(), tx, uuid.MustParse(id))
		if err != nil {
			//TODO: give client more inforamtion instead of an timeout
			slog.Error("Error getting shopping cart", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return err
		}

		transaction := cart.GenerateTransaction()

		a.UserRepo.CreateTransaction(r.Context(), tx, transaction)

		cart.Checkout(transaction.Id)

		a.Repo.SaveShoppingCart(r.Context(), tx, cart)

		a.UserRepo.CreateTransaction(r.Context(), tx, transaction)

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
