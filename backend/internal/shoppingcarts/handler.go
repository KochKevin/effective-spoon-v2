package shoppingcarts

import (
	"context"
	"errors"
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

func (a *Api) ToDto(cart ShoppingCart, freeAmountUsed int, totalFreeAmount int) (dto shoppingcartsapi.ShoppingCart) {

	var lineItems []shoppingcartsapi.LineItem

	for _, item := range cart.LineItems {

		lineItems = append(lineItems, shoppingcartsapi.LineItem{
			Amount:      item.Amount.GetTotalAmount(),
			Price:       float32(item.GetPrice().GetAsEuro()),
			ProductId:   item.Product.Id.String(),
			ProductName: item.Product.Name,
		})
	}

	return shoppingcartsapi.ShoppingCart{
		Id:              cart.Id.String(),
		FullPrice:       float32(cart.GetFullPrice().GetAsEuro()),
		LineItems:       lineItems,
		UserId:          cart.UserId.String(),
		TransactionId:   cart.TransactionId.UUID.String(),
		Status:          string(cart.Status),
		FreeAmountUsed:  freeAmountUsed,
		TotalFreeAmount: totalFreeAmount,
		UseEvent:        cart.UseEvent,
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

	var dto shoppingcartsapi.ShoppingCart

	if cart.UseEvent {

		event, err := a.EventService.GetEvent(r.Context(), cart.EventId)
		if err != nil {
			slog.Error("getting event", "error", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		eventUsage, err := a.EventService.GetEventUsage(r.Context(), userID, cart.EventId)
		if err != nil {
			slog.Error("getting event usage", "error", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		dto = a.ToDto(cart, eventUsage.AmountUsed, event.AmountFreeProductsPerUser)

	} else {
		dto = a.ToDto(cart, 0, 0)
	}

	render.JSON(w, r, dto)
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

	var dto shoppingcartsapi.ShoppingCart

	if cart.UseEvent {

		event, err := a.EventService.GetEvent(r.Context(), cart.EventId)
		if err != nil {
			slog.Error("getting event", "error", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		eventUsage, err := a.EventService.GetEventUsage(r.Context(), userID, cart.EventId)
		if err != nil {
			slog.Error("getting event usage", "error", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		dto = a.ToDto(cart, eventUsage.AmountUsed, event.AmountFreeProductsPerUser)

	} else {
		dto = a.ToDto(cart, 0, 0)
	}

	render.JSON(w, r, dto)
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

	var dto shoppingcartsapi.ShoppingCart

	if cart.UseEvent {

		event, err := a.EventService.GetEvent(r.Context(), cart.EventId)
		if err != nil {
			slog.Error("getting event", "error", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		eventUsage, err := a.EventService.GetEventUsage(r.Context(), userID, cart.EventId)
		if err != nil {
			slog.Error("getting event usage", "error", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		dto = a.ToDto(cart, eventUsage.AmountUsed, event.AmountFreeProductsPerUser)

	} else {
		dto = a.ToDto(cart, 0, 0)
	}

	render.JSON(w, r, dto)

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

	var dto shoppingcartsapi.ShoppingCart

	if cart.UseEvent {

		event, err := a.EventService.GetEvent(r.Context(), cart.EventId)
		if err != nil {
			slog.Error("getting event", "error", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		eventUsage, err := a.EventService.GetEventUsage(r.Context(), userID, cart.EventId)
		if err != nil {
			slog.Error("getting event usage", "error", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		dto = a.ToDto(cart, eventUsage.AmountUsed, event.AmountFreeProductsPerUser)

	} else {
		dto = a.ToDto(cart, 0, 0)
	}

	render.JSON(w, r, dto)
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
	var dto shoppingcartsapi.ShoppingCart

	if cart.UseEvent {

		event, err := a.EventService.GetEvent(r.Context(), cart.EventId)
		if err != nil {
			slog.Error("getting event", "error", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		eventUsage, err := a.EventService.GetEventUsage(r.Context(), userID, cart.EventId)
		if err != nil {
			slog.Error("getting event usage", "error", err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		dto = a.ToDto(cart, eventUsage.AmountUsed, event.AmountFreeProductsPerUser)

	} else {
		dto = a.ToDto(cart, 0, 0)
	}

	render.JSON(w, r, dto)
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
