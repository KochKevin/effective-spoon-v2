package shoppingcarts

import (
	"context"
	"errors"
	"net/http"

	"log/slog"

	shoppingcartsapi "github.com/KochKevin/effective-spoon-v2/internal/shoppingcarts/generated"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type ShoppingCartService interface {
	CreateCurrentShoppingCart(ctx context.Context, userId uuid.UUID) (cart ShoppingCart, err error)
	CheckoutCurrentShoppingCart(ctx context.Context, userId uuid.UUID) (cart ShoppingCart, err error)
	DecreaseProductOfCurrentShoppingCart(ctx context.Context, userId uuid.UUID, productId uuid.UUID) (cart ShoppingCart, err error)
	IncreaseProductOfCurrentShoppingCart(ctx context.Context, userId uuid.UUID, productId uuid.UUID) (cart ShoppingCart, err error)
}

type Api struct {
	Service ShoppingCartService
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

	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		err := errors.New("cannot get user_id from context")
		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	cart, err := a.Service.CreateCurrentShoppingCart(r.Context(), userID)
	if err != nil {
		slog.Error("error while Creating new shopping cart", "error", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, a.ToDto(cart))
}

// Check out of the current shopping cart
// (POST /shopping-carts/current/checkout)
func (a *Api) PostShoppingCartsCurrentCheckout(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		err := errors.New("cannot get user_id from context")
		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	cart, err := a.Service.CheckoutCurrentShoppingCart(r.Context(), userID)
	if err != nil {
		slog.Error("error while checking out current shopping cart", "error", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, a.ToDto(cart))
}

// Remove product from the current shopping cart
// (POST /shopping-carts/current/decrease)
func (a *Api) PostShoppingCartsCurrentDecrease(w http.ResponseWriter, r *http.Request, params shoppingcartsapi.PostShoppingCartsCurrentDecreaseParams) {

	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		err := errors.New("cannot get user_id from context")
		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	cart, err := a.Service.DecreaseProductOfCurrentShoppingCart(r.Context(), userID, uuid.MustParse(params.ProductID))
	if err != nil {
		slog.Error("error while checking out current shopping cart", "error", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, a.ToDto(cart))

}

// Add product to the current shopping cart
// (POST /shopping-carts/current/increase)
func (a *Api) PostShoppingCartsCurrentIncrease(w http.ResponseWriter, r *http.Request, params shoppingcartsapi.PostShoppingCartsCurrentIncreaseParams) {

	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		err := errors.New("cannot get user_id from context")
		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	cart, err := a.Service.IncreaseProductOfCurrentShoppingCart(r.Context(), userID, uuid.MustParse(params.ProductID))
	if err != nil {
		slog.Error("error while checking out current shopping cart", "error", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, a.ToDto(cart))
}
