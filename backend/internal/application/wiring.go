package application

import (
	"context"
	"fmt"

	"github.com/KochKevin/effective-spoon-v2/internal/auth/authcache"
	authservice "github.com/KochKevin/effective-spoon-v2/internal/auth/service"
	"github.com/KochKevin/effective-spoon-v2/internal/chargements/chargementcache"
	chargementservice "github.com/KochKevin/effective-spoon-v2/internal/chargements/service"
	chargementssqlite "github.com/KochKevin/effective-spoon-v2/internal/chargements/sqlite"
	"github.com/KochKevin/effective-spoon-v2/internal/events/eventscache"
	eventslifecycle "github.com/KochKevin/effective-spoon-v2/internal/events/lifecycle"
	eventsservice "github.com/KochKevin/effective-spoon-v2/internal/events/service"
	eventssqlite "github.com/KochKevin/effective-spoon-v2/internal/events/sqlite"
	eventsusage "github.com/KochKevin/effective-spoon-v2/internal/events/usage"
	sqlc "github.com/KochKevin/effective-spoon-v2/internal/infrastructure/sqlite/generated"
	inputservice "github.com/KochKevin/effective-spoon-v2/internal/input/service"
	productssqlite "github.com/KochKevin/effective-spoon-v2/internal/products/sqlite"
	pushservice "github.com/KochKevin/effective-spoon-v2/internal/push/service"
	shoppingcartservice "github.com/KochKevin/effective-spoon-v2/internal/shoppingcarts/service"
	"github.com/KochKevin/effective-spoon-v2/internal/shoppingcarts/shoppingcartcache"
	shoppingcartssqlite "github.com/KochKevin/effective-spoon-v2/internal/shoppingcarts/sqlite"
	userservice "github.com/KochKevin/effective-spoon-v2/internal/users/service"
	userssqlite "github.com/KochKevin/effective-spoon-v2/internal/users/sqlite"
	"github.com/stripe/stripe-go/v86"
)

type Repositories struct {
	AuthCache         *authcache.Cache
	Users             *userssqlite.Repo
	Events            *eventssqlite.Repo
	EventCache        *eventscache.Cache
	Products          *productssqlite.Repo
	ShoppingCarts     *shoppingcartssqlite.Repo
	ShoppingCartCache *shoppingcartcache.Cache
	Chargements       *chargementssqlite.Repo
	ChargementCache   *chargementcache.Cache
}

func newRepos(queries *sqlc.Queries) Repositories {
	return Repositories{
		AuthCache:         authcache.New(),
		Users:             &userssqlite.Repo{Queries: *queries},
		Events:            &eventssqlite.Repo{Queries: *queries},
		EventCache:        eventscache.New(),
		Products:          &productssqlite.Repo{Queries: *queries},
		ShoppingCarts:     &shoppingcartssqlite.Repo{Queries: *queries},
		ShoppingCartCache: shoppingcartcache.New(),
		Chargements:       &chargementssqlite.Repo{Queries: *queries},
		ChargementCache:   chargementcache.New(),
	}
}

type Services struct {
	Auth           *authservice.AuthService
	User           *userservice.Service
	Push           *pushservice.PushService
	Event          *eventsservice.Service
	EventUsage     *eventsusage.UsageService
	EventLifecycle *eventslifecycle.Service
	Shoppingcart   *shoppingcartservice.ShoppingCartService
	Input          *inputservice.InputService
	Chargement     *chargementservice.Service
}

func (a *Application) newServices(ctx context.Context) (Services, error) {
	authService := &authservice.AuthService{
		Repo: a.Repos.AuthCache,
	}

	pushService := pushservice.New()

	userService := &userservice.Service{
		Repo: a.Repos.Users,
		Txm:  a.Txm,
	}

	eventsService := eventsservice.NewService(
		a.Txm,
		a.Repos.Events,
		a.Repos.EventCache,
		userService,
	)

	eventsUsageService := eventsusage.NewService(
		a.Txm,
		a.Repos.Events,
	)

	shoppingCartService := shoppingcartservice.NewService(
		a.Repos.ShoppingCarts,
		a.Repos.Products,
		a.Repos.Users,
		a.Repos.ShoppingCartCache,
		eventsService,
		eventsUsageService,
		a.Txm,
	)

	eventLifecycle, err := eventslifecycle.NewService(
		ctx,
		a.Txm,
		eventsService,
		shoppingCartService,
		pushService,
	)

	if err != nil {
		return Services{}, fmt.Errorf("in wiring event lifecycle service: %w", err)
	}

	inputService := inputservice.NewService(
		a.Repos.Products,
		shoppingCartService,
		authService,
		a.Repos.Users,
		pushService,
		a.Txm,
	)

	chargementService, err := chargementservice.NewService(
		ctx,
		a.Txm,
		a.Repos.Chargements,
		a.Repos.ChargementCache,
		a.Repos.Users,
		stripe.NewClient(a.cfg.StripeClientKey),
		pushService,
	)

	if err != nil {
		return Services{}, fmt.Errorf("in wiring chargement service: %w", err)
	}

	return Services{
		Auth:           authService,
		User:           userService,
		Push:           pushService,
		Event:          eventsService,
		EventUsage:     eventsUsageService,
		EventLifecycle: eventLifecycle,
		Shoppingcart:   shoppingCartService,
		Input:          inputService,
		Chargement:     chargementService,
	}, nil
}
