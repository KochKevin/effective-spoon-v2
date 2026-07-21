package inputservice

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"strings"

	authservice "github.com/KochKevin/effective-spoon-v2/internal/auth/service"
	"github.com/KochKevin/effective-spoon-v2/internal/infrastructure"
	"github.com/KochKevin/effective-spoon-v2/internal/products"
	"github.com/KochKevin/effective-spoon-v2/internal/shoppingcarts"
	shoppingcartservice "github.com/KochKevin/effective-spoon-v2/internal/shoppingcarts/service"
	"github.com/google/uuid"
)

type ProductRepo interface {
	GetProductByCode(ctx context.Context, tx *sql.Tx, code string) (products.Product, error)
}

type AuthService interface {
	GetCurrentUserId() (uuid.UUID, error)
}

type ShoppingCartService interface {
	GetCurrentShoppingCartId() (uuid.UUID, error)
	IncreaseProductOfCurrentShoppingCartTx(ctx context.Context, tx *sql.Tx, userId uuid.UUID, productId uuid.UUID) (cart shoppingcarts.ShoppingCart, err error)
}

type InputService struct {
	ProductRepo         ProductRepo
	ShoppingCartService ShoppingCartService
	AuthService         AuthService
	Txm                 infrastructure.TxManager
}

func New(productRepo ProductRepo, shoppingCartService ShoppingCartService, authService AuthService, txm infrastructure.TxManager) *InputService {
	return &InputService{
		ProductRepo:         productRepo,
		ShoppingCartService: shoppingCartService,
		AuthService:         authService,
		Txm:                 txm,
	}
}

func (i *InputService) EnterBarcode(input string) error {

	ctx := context.Background()

	err := i.Txm.WithTx(context.Background(), func(tx *sql.Tx) error {

		userId, err := i.AuthService.GetCurrentUserId()
		if errors.Is(err, authservice.NoCurrentUser) {
			slog.Error("error: to add an product via the barcode scanner, an user must be logged in", "error:", err)
			return err
		}

		_, err = i.ShoppingCartService.GetCurrentShoppingCartId()
		if errors.Is(err, shoppingcartservice.NoCurrentShoppingCart) {
			slog.Error("error: to add an product to the current shopping cart, an current shopping cart must exist", "error:", err)
			return err
		}

		//Sanitze Input
		input = strings.TrimSpace(input)

		//Load Product
		product, err := i.ProductRepo.GetProductByCode(ctx, tx, input)
		if err != nil {
			slog.Error("error when getting an product by its code", "error:", err)
			return err
		}

		//Add Product to Cart
		i.ShoppingCartService.IncreaseProductOfCurrentShoppingCartTx(ctx, tx, userId, product.Id)

		return nil
	})
	if err != nil {
		slog.Error("error when using input system when entering barcode", "error:", err)
		return err
	}

	return nil
}

func (i *InputService) EnterRfid(input string) error {

	err := i.Txm.WithTx(context.Background(), func(tx *sql.Tx) error {

		return nil
	})
	if err != nil {
		slog.Error("error when using input system when entering rfid", "error:", err)
		return err
	}

	return nil
}
