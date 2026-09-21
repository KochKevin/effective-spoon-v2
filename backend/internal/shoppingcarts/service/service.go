package shoppingcartservice

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/KochKevin/effective-spoon-v2/internal/events"
	"github.com/KochKevin/effective-spoon-v2/internal/infrastructure"
	"github.com/KochKevin/effective-spoon-v2/internal/products"
	"github.com/KochKevin/effective-spoon-v2/internal/shoppingcarts"
	"github.com/KochKevin/effective-spoon-v2/internal/users"
	"github.com/google/uuid"
)

type Repo interface {
	CreateShoppingCart(ctx context.Context, tx *sql.Tx, cart shoppingcarts.ShoppingCart) (shoppingcarts.ShoppingCart, error)
	GetShoppingCart(ctx context.Context, tx *sql.Tx, id uuid.UUID) (shoppingcarts.ShoppingCart, error)
	SaveShoppingCart(ctx context.Context, tx *sql.Tx, cart shoppingcarts.ShoppingCart) error
	DeleteShoppingCart(ctx context.Context, tx *sql.Tx, id uuid.UUID) error
}

type ProductRepo interface {
	GetProduct(ctx context.Context, tx *sql.Tx, id uuid.UUID) (products.Product, error)
}

type UserRepo interface {
	CreateTransaction(ctx context.Context, tx *sql.Tx, transaction users.Transaction) (users.Transaction, error)
}

type EventService interface {
	GetEventUsage(ctx context.Context, userId uuid.UUID, eventId uuid.UUID) (eventUsage events.EventUsage, err error)
	GetEvent(ctx context.Context, eventId uuid.UUID) (event events.Event, err error)
	GetCurrentEvent(ctx context.Context) (event events.Event, err error)
}

type ShoppingCartCache interface {
	SetCurrentCartId(cartId uuid.UUID)
	GetCurrentCartId() (cartId uuid.UUID)
	ClearCurrentCartId()
}

type ShoppingCartService struct {
	Repo              Repo
	ProductRepo       ProductRepo
	UserRepo          UserRepo
	EventService      EventService
	ShoppingCartCache ShoppingCartCache
	Txm               infrastructure.TxManager
}

func New(repo Repo, productRepo ProductRepo, userRepo UserRepo, shoppingCartCache ShoppingCartCache, eventSerice EventService, txm infrastructure.TxManager) *ShoppingCartService {
	return &ShoppingCartService{
		Repo:              repo,
		ProductRepo:       productRepo,
		UserRepo:          userRepo,
		ShoppingCartCache: shoppingCartCache,
		EventService:      eventSerice,
		Txm:               txm,
	}
}

func (s *ShoppingCartService) checkIfShoppingCartCanBeUsed(cart shoppingcarts.ShoppingCart, userId uuid.UUID) error {

	//Check if cart owner and current user are the same
	if cart.UserId != userId {
		return errors.New("shopping cart does not belong to the logged in user")
	}

	//Check if cart is active
	if cart.Status != shoppingcarts.ShoppingCartActive {
		return errors.New("shopping cart is not active, can not check out")
	}

	return nil
}

// Crate the current Shopping cart
func (s *ShoppingCartService) CreateCurrentShoppingCart(ctx context.Context, userId uuid.UUID) (cart shoppingcarts.ShoppingCart, err error) {

	// Connect an shopping cart to an event only when an event is active to the creation time of the shopping cart
	currentEvent, err := s.EventService.GetCurrentEvent(ctx)

	slog.Debug("Create Current Shopping Cart - Get Current Event error: ", "error", err)

	if errors.Is(err, events.NoCurrentEventErr) {

		cart = shoppingcarts.NewShoppingCart(userId, uuid.Nil, false)
	} else {
		cart = shoppingcarts.NewShoppingCart(userId, currentEvent.Id, true)
	}

	//Only break if an real error appears
	if err != nil && !errors.Is(err, events.NoCurrentEventErr) {

		slog.Error("getting current event", "error:", err)
		return shoppingcarts.ShoppingCart{}, err
	}

	slog.Debug("Create Current Shopping Cart", "cart", cart)

	err = s.Txm.WithTx(ctx, func(tx *sql.Tx) error {
		_, err = s.Repo.CreateShoppingCart(ctx, tx, cart)
		if err != nil {
			return fmt.Errorf("in creating shopping cart: %w", err)
		}

		//Ignore if an current cart is already set, replace it!
		s.ShoppingCartCache.SetCurrentCartId(cart.Id)

		return nil

	})

	if err != nil {

		return shoppingcarts.ShoppingCart{}, fmt.Errorf("in transaction: %w")
	}

	return cart, nil

}

func (s *ShoppingCartService) CheckoutCurrentShoppingCart(ctx context.Context, userId uuid.UUID) (cart shoppingcarts.ShoppingCart, err error) {

	err = s.Txm.WithTx(ctx, func(tx *sql.Tx) error {

		//Get Current Cart id
		cartId, err := s.GetCurrentShoppingCartId()
		if err != nil {
			slog.Error("error can not get current shopping cart id", "error:", err)
			return err
		}

		cart, err = s.Repo.GetShoppingCart(ctx, tx, cartId)
		if err != nil {
			slog.Error("Error getting shopping cart", "error:", err)
			return err
		}

		err = s.checkIfShoppingCartCanBeUsed(cart, userId)
		if err != nil {
			slog.Error("Error cart can not be used for checkout", "error:", err)
			return err
		}

		transaction := cart.GenerateTransaction()

		transaction, err = s.UserRepo.CreateTransaction(ctx, tx, transaction)
		if err != nil {
			slog.Error("Error creating transaction", "error", err)
			return err
		}

		cart.Checkout(transaction.Id)

		err = s.Repo.SaveShoppingCart(ctx, tx, cart)
		if err != nil {
			slog.Error("Error saving shoppingcart", "error", err)
			return err
		}

		s.ShoppingCartCache.ClearCurrentCartId()

		return nil

	})

	if err != nil {
		slog.Error("Error in /shopping carts transaction", "error:", err)
		return shoppingcarts.ShoppingCart{}, err
	}

	return cart, nil
}
func (s *ShoppingCartService) DecreaseProductOfCurrentShoppingCart(ctx context.Context, userId uuid.UUID, productId uuid.UUID) (cart shoppingcarts.ShoppingCart, err error) {
	slog.Debug("Decrease Product by id")

	err = s.Txm.WithTx(ctx, func(tx *sql.Tx) error {

		//Get Current Cart id
		cartId, err := s.GetCurrentShoppingCartId()
		if err != nil {
			slog.Error("error can not get current shopping cart id", "error:", err)
			return err
		}

		cart, err = s.Repo.GetShoppingCart(ctx, tx, cartId)
		if err != nil {
			slog.Error("Error getting shopping cart", "error:", err)
			return err
		}

		err = s.checkIfShoppingCartCanBeUsed(cart, userId)
		if err != nil {
			slog.Error("Error cart can not be used for checkout", "error:", err)
			return err
		}

		slog.Debug("Before decrease: ", "cart", cart)

		cart.DecreaseProductAmount(productId)

		err = s.Repo.SaveShoppingCart(ctx, tx, cart)
		if err != nil {
			slog.Error("Error getting shopping cart", "error:", err)
			return err
		}
		return nil
	})

	if err != nil {
		slog.Error("Error in /shoppingcart decrease transaction", "error:", err)
		return shoppingcarts.ShoppingCart{}, err
	}

	return cart, nil

}

func (s *ShoppingCartService) IncreaseProductOfCurrentShoppingCart(ctx context.Context, userId uuid.UUID, productId uuid.UUID) (cart shoppingcarts.ShoppingCart, err error) {

	err = s.Txm.WithTx(ctx, func(tx *sql.Tx) error {

		cart, err = s.IncreaseProductOfCurrentShoppingCartTx(ctx, tx, userId, productId)
		return err
	})

	if err != nil {
		slog.Error("Error in shoppingcart decrease transaction", "error:", err)
		return shoppingcarts.ShoppingCart{}, err
	}

	return cart, err
}

func (s *ShoppingCartService) IncreaseProductOfCurrentShoppingCartTx(ctx context.Context, tx *sql.Tx, userId uuid.UUID, productId uuid.UUID) (cart shoppingcarts.ShoppingCart, err error) {

	slog.Debug("Increase Product by id")

	//Get Current Cart id
	cartId, err := s.GetCurrentShoppingCartId()
	if err != nil {
		slog.Error("error can not get current shopping cart id", "error:", err)
		return shoppingcarts.ShoppingCart{}, err
	}

	cart, err = s.Repo.GetShoppingCart(ctx, tx, cartId)
	if err != nil {
		//fmt.Errorf("error getting shopping cart: %w", err)
		//slog.Error("error getting shopping cart", "error:", err)
		return shoppingcarts.ShoppingCart{}, fmt.Errorf("error getting shopping cart: %w", err)
	}

	err = s.checkIfShoppingCartCanBeUsed(cart, userId)
	if err != nil {
		slog.Error("error cart can not be used for checkout", "error:", err)
		return shoppingcarts.ShoppingCart{}, err
	}

	slog.Debug("Before increase: ", "cart", cart)

	product, err := s.ProductRepo.GetProduct(ctx, tx, productId)
	if err != nil {
		slog.Error("error cart can not load product which should be added to cart", "error:", err)
		return shoppingcarts.ShoppingCart{}, err
	}

	if cart.UseEvent {

		//Get the amount of free products the user already used
		eventUsage, err := s.EventService.GetEventUsage(ctx, userId, cart.EventId)
		if err != nil {
			slog.Error("getting event usage", "error:", err)
			return shoppingcarts.ShoppingCart{}, err
		}

		//Get the event which is connected with the shopping cart
		event, err := s.EventService.GetEvent(ctx, cart.EventId)
		if err != nil {
			slog.Error("getting event", "error:", err)
			return shoppingcarts.ShoppingCart{}, err
		}

		//Use free products first before adding products which the user needs to pay for
		if eventUsage.AmountUsed >= event.AmountFreeProductsPerUser {
			cart.IncreaseProductAmount(product, false)
		} else {
			cart.IncreaseProductAmount(product, true)
		}

	} else {
		cart.IncreaseProductAmount(product, false)
	}

	err = s.Repo.SaveShoppingCart(ctx, tx, cart)
	if err != nil {
		slog.Error("Error getting shopping cart", "error:", err)
		return shoppingcarts.ShoppingCart{}, err
	}

	return cart, nil

}

var NoCurrentShoppingCart = errors.New("error: no current shopping cart is setted to be getted. Set an current shopping cart first")

func (s *ShoppingCartService) GetCurrentShoppingCartId() (uuid.UUID, error) {
	cartId := s.ShoppingCartCache.GetCurrentCartId()

	if cartId == uuid.Nil {
		return uuid.Nil, NoCurrentShoppingCart
	}

	return cartId, nil
}

func (s *ShoppingCartService) GetCurrentShoppingCart(ctx context.Context, userId uuid.UUID) (cart shoppingcarts.ShoppingCart, err error) {

	err = s.Txm.WithTx(ctx, func(tx *sql.Tx) error {

		//Get Current Cart id
		cartId, err := s.GetCurrentShoppingCartId()
		if err != nil {
			slog.Error("error can not get current shopping cart id", "error:", err)
			return err
		}

		cart, err = s.Repo.GetShoppingCart(ctx, tx, cartId)
		if err != nil {
			slog.Error("error getting shopping cart", "error:", err)
			return err
		}

		return err
	})

	if err != nil {
		slog.Error("Error in shoppingcart decrease transaction", "error:", err)
		return shoppingcarts.ShoppingCart{}, err
	}

	return cart, err

}

func (s *ShoppingCartService) CancelCurrentShoppingCart(ctx context.Context, userId uuid.UUID) (err error) {

	err = s.Txm.WithTx(ctx, func(tx *sql.Tx) error {

		//Get Current Cart id
		cartId, err := s.GetCurrentShoppingCartId()
		if err != nil {

			return fmt.Errorf("can not get current shopping cart id: %w", err)
		}

		cart, err := s.Repo.GetShoppingCart(ctx, tx, cartId)
		if err != nil {
			return fmt.Errorf("getting shopping cart: %w", err)
		}

		err = s.checkIfShoppingCartCanBeUsed(cart, userId)
		if err != nil {
			return fmt.Errorf("cart can not be used: %w", err)
		}

		err = s.Repo.DeleteShoppingCart(ctx, tx, cartId)
		if err != nil {
			return fmt.Errorf("in deleting: %w", err)
		}

		s.ShoppingCartCache.ClearCurrentCartId()

		return nil

	})

	if err != nil {
		return fmt.Errorf("in transaction: %w", err)
	}

	return nil

}
