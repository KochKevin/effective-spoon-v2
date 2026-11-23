package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/KochKevin/effective-spoon-v2/internal/application"
	"github.com/KochKevin/effective-spoon-v2/internal/config"
	"github.com/KochKevin/effective-spoon-v2/internal/infrastructure"
	"github.com/google/uuid"
	"github.com/pressly/goose/v3"

	_ "modernc.org/sqlite"
)

type AuthService interface {
	GetCurrentUserId() uuid.UUID
}

func main() {

	ctx := context.Background()

	//Setup Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	slog.Info("Backend Api started")

	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	//Create database folder

	dbFolderPath := filepath.Join(".", "db")
	err = os.MkdirAll(dbFolderPath, os.ModePerm)
	if err != nil {
		slog.Error("can not create folder for data.db", "err", err)
		return
	}

	//Do Database migrations. Open sqlite with WAL mode for writing and reading

	sqliteConnectionString := "file:" + dbFolderPath + "/data.db?_pragma=journal_mode=WAL&_pragma=busy_timeout=5000"

	db, err := goose.OpenDBWithDriver("sqlite", sqliteConnectionString)
	if err != nil {
		slog.Error("can not open data.db", "err", err)
		return
	}

	defer func() { _ = db.Close() }()

	migrationsFS, err := fs.Sub(infrastructure.DBMigrations, "db/migrations")
	if err != nil {
		slog.Error("failed to create sub filesystem for migrations", "err", err)
		return
	}

	gooseProvider, err := goose.NewProvider(goose.DialectSQLite3, db, migrationsFS, goose.WithSlog(logger))
	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = gooseProvider.Up(ctx)

	if err != nil {
		slog.Error("Error migrating database", "error", err)
		return
	}

	app, err := application.New(ctx, db, cfg)
	if err != nil {
		slog.Error("Error in app", "error", fmt.Errorf("in building app: %w", err))
		return
	}

	app.Run()

	/*

		//Router Setup
		r := chi.NewRouter()

		//Middlewear
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://localhost:*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
		}))

		//Serve everything under /api

		r.Route("/api", func(apiRouter chi.Router) {

			authService := &authservice.AuthService{
				Repo: authcache.New(),
			}

			pushService := pushservice.New()

			//The auth api should only be used for testing purposes
			authapi.HandlerFromMux(&auth.Api{
				UserRepo: &userssqlite.Repo{
					Queries: *sqlc.New(db),
				},
				AuthService: authService,
				Txm:         *infrastructure.NewTxManager(db),
			}, apiRouter)

			//User Service

			userService := userservice.Service{
				Repo: &userssqlite.Repo{
					Queries: *sqlc.New(db),
				},
				Txm: *infrastructure.NewTxManager(db),
			}

			//Events
			eventsService := eventsservice.NewService(
				*infrastructure.NewTxManager(db),
				&eventssqlite.Repo{
					Queries: *sqlc.New(db),
				},
				eventscache.New(),
				new(userService),
			)

			//Event Usage
			eventsUsageService := eventsusage.NewService(
				*infrastructure.NewTxManager(db),
				&eventssqlite.Repo{
					Queries: *sqlc.New(db),
				},
			)

			//Shoppingcart
			shoppingCartService := shoppingcartservice.NewService(
				&shoppingcartssqlite.Repo{
					Queries: *sqlc.New(db),
				},
				&productssqlite.Repo{
					Queries: *sqlc.New(db),
				},
				&userssqlite.Repo{
					Queries: *sqlc.New(db),
				},
				shoppingcartcache.New(),
				&eventsService,
				&eventsUsageService,
				*infrastructure.NewTxManager(db),
			)

			//Event Lifecycle

			eventLifecycle := eventslifecycle.NewService(
				*infrastructure.NewTxManager(db),
				&eventsService,
				shoppingCartService,
				pushService,
			)

			err := eventLifecycle.Initalize(context.Background())
			if err != nil {
				slog.Error("error while Initlizing Event Lifecycle Service", "err", err)
			}

			inputapi.HandlerFromMux(
				&input.Api{
					InputService: inputservice.NewService(
						&productssqlite.Repo{
							Queries: *sqlc.New(db),
						},
						shoppingCartService,
						authService,
						&userssqlite.Repo{
							Queries: *sqlc.New(db),
						},
						pushService,
						*infrastructure.NewTxManager(db),
					),
				}, apiRouter)

			//Server Sent Events, Push Api endpoint
			pushapi.HandlerFromMux(&push.Api{PushService: pushService}, apiRouter)

			apiRouter.Group(func(protectedRouter chi.Router) {

				//GetLoggedInUser Middlewear
				protectedRouter.Use(func(next http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

						// Create/Get test user id
						//userId := uuid.Nil

						userId, err := authService.GetCurrentUserId()
						if err != nil {
							slog.Error("error: a user needs to be logged in to access this endpoint", "error", err)
							http.Error(w, "Internal Server Error - a user need to be logged in", http.StatusUnauthorized)
							return
						}

						slog.Debug("auth middlewear", "user_id", userId)

						r = r.WithContext(authservice.WithUser(r.Context(), userId))

						next.ServeHTTP(w, r)
					})
				})

				protectedRouter.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(`{"status": "pong"}`))
				})

				//Routes
				productsapi.HandlerFromMux(
					&products.Api{
						Repo: &productssqlite.Repo{
							Queries: *sqlc.New(db),
						},
						Txm: *infrastructure.NewTxManager(db),
					}, protectedRouter)

				shoppingcartssapi.HandlerFromMux(
					&shoppingcarts.Api{
						Service: shoppingCartService,
					}, protectedRouter)

				//Users
				userssapi.HandlerFromMux(
					&users.Api{
						Service: &userService,
					}, protectedRouter)

				//Chargements
				chargementService := chargementservice.Service{
					Repo: &chargementssqlite.Repo{
						Queries: *sqlc.New(db),
					},
					BalanceChargementIntentCache: chargementcache.New(),
					UserRepo: &userssqlite.Repo{
						Queries: *sqlc.New(db),
					},
					Txm:          *infrastructure.NewTxManager(db),
					StripeClient: stripe.NewClient(cfg.StripeClientKey),
					PushService:  pushService,
				}

				chargementService.Ticker(context.Background())

				chargementsapi.HandlerFromMux(&chargements.Api{
					Service: &chargementService,
				}, protectedRouter)

				//Events handler
				eventsapi.HandlerFromMux(&events.Api{
					Service:     &eventsService,
					UserService: &userService,
					Usage:       &eventsUsageService,
					Lifecycle:   new(eventLifecycle),
				}, protectedRouter)

			})

		})

		//Frontend
		server.ServeFrontend(r)

		//Serve
		err = http.ListenAndServe(":8080", r)

		slog.Error("Error serving api", "error", err)
	*/

}
