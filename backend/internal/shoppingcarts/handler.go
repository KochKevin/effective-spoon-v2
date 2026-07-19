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

type ShoppingCartCache interface {
	SetCurrentCartId(cartId uuid.UUID)
	GetCurrentCartId() (cartId uuid.UUID)
	ClearCurrentCartId()
}

type Api struct {
	Repo              Repo
	ProductRepo       ProductRepo
	UserRepo          UserRepo
	ShoppingCartCache ShoppingCartCache
	Txm               infrastructure.TxManager
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
		Id:            cart.Id.String(),
		FullPrice:     float32(cart.GetFullPrice().GetAsEuro()),
		LineItems:     lineItems,
		UserId:        cart.UserId.String(),
		TransactionId: cart.TransactionId.UUID.String(),
		Status:        string(cart.Status),
	}

}

// Create a new shopping cart and set it as the current cart
// (POST /shopping-carts/current)
func (a *Api) PostShoppingCartsCurrent(w http.ResponseWriter, r *http.Request) {
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
			slog.Error("Error creating shopping cart", "error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return err
		}

		//Ignore if an current cart is already set, replace it!
		a.ShoppingCartCache.SetCurrentCartId(cart.Id)

		dto = a.ToDto(cart)

		return nil

	})

	if err != nil {
		slog.Error("Error in /products transaction", "error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, dto)
}

// Check out of the current shopping cart
// (POST /shopping-carts/current/checkout)
func (a *Api) PostShoppingCartsCurrentCheckout(w http.ResponseWriter, r *http.Request) {
	var dto shoppingcartsapi.ShoppingCart

	err := a.Txm.WithTx(context.Background(), func(tx *sql.Tx) error {

		userID, ok := r.Context().Value("user_id").(uuid.UUID)
		if !ok {
			slog.Error("Can not get user_id from context")
			return errors.New("Can not get user_id from context")
		}

		//Get Current Cart id
		cartId := a.ShoppingCartCache.GetCurrentCartId()

		if cartId == uuid.Nil {
			slog.Error("Error: no current Cart id is set")
			http.Error(w, "Internal Server Error, to check out the current shopping cart, a current shopping cart must be set", http.StatusInternalServerError)
			return errors.New("Error: no current Cart id is set")
		}

		cart, err := a.Repo.GetShoppingCart(r.Context(), tx, cartId)
		if err != nil {
			//TODO: give client more inforamtion instead of an timeout
			slog.Error("Error getting shopping cart", "error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return err
		}

		//Check if cart owner and current user are the same
		if cart.UserId != userID {
			slog.Error("Error: shopping cart does not belong to the logged in user")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return errors.New("shopping cart does not belong to the logged in user")
		}

		//Check if cart is active
		if cart.Status != ShoppingCartActive {
			slog.Error("Error shopping cart is not active, can not check out")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return errors.New("shopping cart is not active")
		}

		transaction := cart.GenerateTransaction()

		transaction, err = a.UserRepo.CreateTransaction(r.Context(), tx, transaction)
		if err != nil {
			slog.Error("Error creating transaction", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return err
		}

		cart.Checkout(transaction.Id)

		err = a.Repo.SaveShoppingCart(r.Context(), tx, cart)
		if err != nil {
			slog.Error("Error saving shoppingcart", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return err
		}

		//Clear the cart id from state/cache
		a.ShoppingCartCache.ClearCurrentCartId()

		dto = a.ToDto(cart)

		return nil

	})

	if err != nil {
		slog.Error("Error in /products transaction", "error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, dto)
}

// Remove product from the current shopping cart
// (POST /shopping-carts/current/decrease)
func (a *Api) PostShoppingCartsCurrentDecrease(w http.ResponseWriter, r *http.Request, params shoppingcartsapi.PostShoppingCartsCurrentDecreaseParams) {
	slog.Debug("Decrease Product by id")

	var dto shoppingcartsapi.ShoppingCart

	err := a.Txm.WithTx(context.Background(), func(tx *sql.Tx) error {

		userID, ok := r.Context().Value("user_id").(uuid.UUID)
		if !ok {
			slog.Error("Can not get user_id from context")
			return errors.New("Can not get user_id from context")
		}

		//Get Current Cart id
		cartId := a.ShoppingCartCache.GetCurrentCartId()

		if cartId == uuid.Nil {
			slog.Error("Error: no current Cart id is set")
			http.Error(w, "Internal Server Error, to check out the current shopping cart, a current shopping cart must be set", http.StatusInternalServerError)
			return errors.New("Error: no current Cart id is set")
		}

		cart, err := a.Repo.GetShoppingCart(r.Context(), tx, cartId)
		if err != nil {
			//TODO: give client more inforamtion instead of an timeout
			slog.Error("Error getting shopping cart", "error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return err
		}

		//Check if cart owner and current user are the same
		if cart.UserId != userID {
			slog.Error("Error: shopping cart does not belong to the logged in user")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return errors.New("shopping cart does not belong to the logged in user")
		}

		//Check if cart is active
		if cart.Status != ShoppingCartActive {
			slog.Error("Error shopping cart is not active, can not check out")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return errors.New("shopping cart is not active")
		}

		slog.Debug("Before decrease: ", "cart", cart)

		cart.DecreaseProductAmount(uuid.MustParse(params.ProductID))

		err = a.Repo.SaveShoppingCart(r.Context(), tx, cart)
		if err != nil {
			//TODO: give client more inforamtion instead of an timeout
			slog.Error("Error getting shopping cart", "error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return err
		}

		dto = a.ToDto(cart)
		slog.Debug("After decrease: ", "dto", dto)

		return nil

	})

	if err != nil {
		slog.Error("Error in /shoppingcart decrease transaction", "error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, dto)
}

// Add product to the current shopping cart
// (POST /shopping-carts/current/increase)
func (a *Api) PostShoppingCartsCurrentIncrease(w http.ResponseWriter, r *http.Request, params shoppingcartsapi.PostShoppingCartsCurrentIncreaseParams) {
	slog.Debug("Increase Product by id")

	var dto shoppingcartsapi.ShoppingCart

	err := a.Txm.WithTx(context.Background(), func(tx *sql.Tx) error {

		userID, ok := r.Context().Value("user_id").(uuid.UUID)
		if !ok {
			slog.Error("Can not get user_id from context")
			return errors.New("Can not get user_id from context")
		}

		//Get Current Cart id
		cartId := a.ShoppingCartCache.GetCurrentCartId()

		if cartId == uuid.Nil {
			slog.Error("Error: no current Cart id is set")
			http.Error(w, "Internal Server Error, to increase the product amount in the current shopping cart, a current shopping cart must be set", http.StatusInternalServerError)
			return errors.New("Error: no current Cart id is set")
		}

		cart, err := a.Repo.GetShoppingCart(r.Context(), tx, cartId)

		slog.Debug("Before increase: ", "error:", cart)

		if err != nil {
			//TODO: give client more inforamtion instead of an timeout
			slog.Error("Error getting shopping cart", "error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return err
		}

		if cart.UserId != userID {
			slog.Error("Error: shopping cart does not belong to the logged in user")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return errors.New("shopping cart does not belong to the logged in user")
		}

		if cart.Status != ShoppingCartActive {
			slog.Error("Error shopping cart is not active, can not increase amounts")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return errors.New("shopping cart is not active")
		}


		productToAdd, err := a.ProductRepo.GetProduct(r.Context(), tx, uuid.MustParse(params.ProductID))
		if err != nil {
			//TODO: give client more inforamtion instead of an timeout
			slog.Error("Error getting product to add to shopping cart", "error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return err
		}

		cart.IncreaseProductAmount(productToAdd)

		err = a.Repo.SaveShoppingCart(r.Context(), tx, cart)
		if err != nil {
			//TODO: give client more inforamtion instead of an timeout
			slog.Error("Error getting shopping cart", "error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return err
		}

		//slog.Debug("After increase: ", "cart", cart)

		dto = a.ToDto(cart)
		slog.Debug("After increase: ", "dto", dto)

		return nil

	})

	if err != nil {
		slog.Error("Error in /products transaction", "error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, dto)
}
