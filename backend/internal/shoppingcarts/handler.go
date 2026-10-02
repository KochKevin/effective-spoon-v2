package shoppingcarts

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"log/slog"

	"github.com/KochKevin/effective-spoon-v2/internal/events"
	shoppingcartsapi "github.com/KochKevin/effective-spoon-v2/internal/shoppingcarts/generated"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type ShoppingCartService interface {
	GetCurrentShoppingCart(ctx context.Context, userId uuid.UUID) (cart ShoppingCart, err error)
	CreateCurrentShoppingCart(ctx context.Context, userId uuid.UUID) (cart ShoppingCart, err error)
	CheckoutCurrentShoppingCart(ctx context.Context, userId uuid.UUID) (cart ShoppingCart, err error)
	CancelCurrentShoppingCart(ctx context.Context, userId uuid.UUID) (err error)
	DecreaseProductOfCurrentShoppingCart(ctx context.Context, userId uuid.UUID, productId uuid.UUID) (cart ShoppingCart, err error)
	IncreaseProductOfCurrentShoppingCart(ctx context.Context, userId uuid.UUID, productId uuid.UUID) (cart ShoppingCart, err error)
	ToggleEventUsageOfCurrentCart(ctx context.Context, userId uuid.UUID, useEvent bool) (cart ShoppingCart, err error)

	GetShoppingCartView(ctx context.Context, cart ShoppingCart) (ShoppingCartView, error)
}

type EventService interface {
	GetEventUsage(ctx context.Context, userId uuid.UUID, eventId uuid.UUID) (eventUsage events.EventUsage, err error)
	GetEvent(ctx context.Context, eventId uuid.UUID) (event events.Event, err error)
}

type Api struct {
	Service      ShoppingCartService
	EventService EventService
}

//a *Api github.com/KochKevin/effective-spoon-v2/internal/shoppingcarts/generated.ServerInterface

func (a *Api) ToDto(cartView ShoppingCartView) (dto shoppingcartsapi.ShoppingCart) {

	var lineItems []shoppingcartsapi.LineItem

	for _, item := range cartView.Cart.LineItems {

		lineItems = append(lineItems, shoppingcartsapi.LineItem{
			Amount:      item.Amount.GetTotalAmount(),
			Price:       float32(item.GetPrice().GetAsEuro()),
			ProductId:   item.Product.Id.String(),
			ProductName: item.Product.Name,
		})

	}

	eventUsageDto := shoppingcartsapi.EventUsage{
		UsedFreeAmount:               cartView.EventUsage.AmountUsed,
		AvailableFreeAmountPerPerson: cartView.FreeProductPerPerson,
	}

	return shoppingcartsapi.ShoppingCart{
		Id:            cartView.Cart.Id.String(),
		FullPrice:     float32(cartView.Cart.GetFullPrice().GetAsEuro()),
		LineItems:     lineItems,
		UserId:        cartView.Cart.UserId.String(),
		TransactionId: cartView.Cart.TransactionId.UUID.String(),
		Status:        string(cartView.Cart.Status),
		UseEvent:      cartView.Cart.UseEvent,
		EventUsage:    eventUsageDto,
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

	cartView, err := a.Service.GetShoppingCartView(r.Context(), cart)
	if err != nil {
		slog.Error("error while getting the shopping cart view", "error", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, a.ToDto(cartView))
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

	cartView, err := a.Service.GetShoppingCartView(r.Context(), cart)
	if err != nil {
		slog.Error("error while getting the shopping cart view", "error", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, a.ToDto(cartView))
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

	cartView, err := a.Service.GetShoppingCartView(r.Context(), cart)
	if err != nil {
		slog.Error("error while getting the shopping cart view", "error", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, a.ToDto(cartView))

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

	cartView, err := a.Service.GetShoppingCartView(r.Context(), cart)
	if err != nil {
		slog.Error("error while getting the shopping cart view", "error", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, a.ToDto(cartView))
}

// GetShoppingCartsCurrent Get the current shopping cart
// (GET /shopping-carts/current)
func (a *Api) GetShoppingCartsCurrent(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		err := errors.New("cannot get user_id from context")
		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	cart, err := a.Service.GetCurrentShoppingCart(r.Context(), userID)
	if err != nil {
		slog.Error("error while getting current shopping cart", "error", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	
	cartView, err := a.Service.GetShoppingCartView(r.Context(), cart)
	if err != nil {
		slog.Error("error while getting the shopping cart view", "error", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, a.ToDto(cartView))
}

// PostShoppingCartsCurrentCancel Cancel and delete the current shopping cart
// (POST /shopping-carts/current/cancel)
func (a *Api) PostShoppingCartsCurrentCancel(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		err := errors.New("cannot get user_id from context")
		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err := a.Service.CancelCurrentShoppingCart(r.Context(), userID)
	if err != nil {
		slog.Error("cancel cart", "error", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, nil)

}

// PutShoppingCartsCurrentUseEvent Set if the shopping cart should use the current event
// (PUT /shopping-carts/current/use-event)
func (a *Api) PutShoppingCartsCurrentUseEvent(w http.ResponseWriter, r *http.Request) {

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("reading request body", "err", err)
	}
	defer r.Body.Close()

	toggelUseEvent := shoppingcartsapi.PutShoppingCartsCurrentUseEventJSONRequestBody{}
	err = json.Unmarshal(bodyBytes, &toggelUseEvent)
	if err != nil {
		slog.Error("unmarshaling create event request body", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		err := errors.New("cannot get user_id from context")
		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	cart, err := a.Service.ToggleEventUsageOfCurrentCart(r.Context(), userID, toggelUseEvent.UseEvent)
	if err != nil {
		slog.Error("shoppingcart use event toggle", "error", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	cartView, err := a.Service.GetShoppingCartView(r.Context(), cart)
	if err != nil {
		slog.Error("error while getting the shopping cart view", "error", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, a.ToDto(cartView))

}
