package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/KochKevin/effective-spoon-v2/internal/auth"
	"github.com/KochKevin/effective-spoon-v2/internal/auth/authcache"
	authapi "github.com/KochKevin/effective-spoon-v2/internal/auth/generated"
	authservice "github.com/KochKevin/effective-spoon-v2/internal/auth/service"
	"github.com/KochKevin/effective-spoon-v2/internal/infrastructure"
	sqlc "github.com/KochKevin/effective-spoon-v2/internal/infrastructure/sqlite/generated"
	"github.com/KochKevin/effective-spoon-v2/internal/input"
	inputapi "github.com/KochKevin/effective-spoon-v2/internal/input/generated"
	inputservice "github.com/KochKevin/effective-spoon-v2/internal/input/service"
	"github.com/KochKevin/effective-spoon-v2/internal/products"
	productsapi "github.com/KochKevin/effective-spoon-v2/internal/products/generated"
	productssqlite "github.com/KochKevin/effective-spoon-v2/internal/products/sqlite"
	"github.com/KochKevin/effective-spoon-v2/internal/push"
	pushapi "github.com/KochKevin/effective-spoon-v2/internal/push/generated"
	pushservice "github.com/KochKevin/effective-spoon-v2/internal/push/service"
	"github.com/KochKevin/effective-spoon-v2/internal/server"
	"github.com/KochKevin/effective-spoon-v2/internal/shoppingcarts"
	shoppingcartssapi "github.com/KochKevin/effective-spoon-v2/internal/shoppingcarts/generated"
	shoppingcartservice "github.com/KochKevin/effective-spoon-v2/internal/shoppingcarts/service"
	"github.com/KochKevin/effective-spoon-v2/internal/shoppingcarts/shoppingcartcache"
	shoppingcartssqlite "github.com/KochKevin/effective-spoon-v2/internal/shoppingcarts/sqlite"
	"github.com/KochKevin/effective-spoon-v2/internal/users"
	userssapi "github.com/KochKevin/effective-spoon-v2/internal/users/generated"
	userssqlite "github.com/KochKevin/effective-spoon-v2/internal/users/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/pressly/goose/v3"

	_ "modernc.org/sqlite"
)

type AuthService interface {
	GetCurrentUserId() uuid.UUID
}

func main() {

	//Setup Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	slog.Info("Backend Api started")

	//Do Database migrations. Open swlite with WAL mode for writing and reading

	sqliteConnectionString := "file:./db/data/data.db?_pragma=journal_mode=WAL&_pragma=busy_timeout=5000"

	db, err := goose.OpenDBWithDriver("sqlite", sqliteConnectionString)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer db.Close()

	err = goose.Up(db, "./db/migrations")

	if err != nil {
		slog.Error("Error migrating database", "error", err)
		return
	}

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

		shoppingCartService := shoppingcartservice.New(
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
			*infrastructure.NewTxManager(db),
		)

		inputapi.HandlerFromMux(
			&input.Api{
				InputService: inputservice.New(
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

					slog.Debug("Auth Middlewear uses ", userId, " UserId")

					r = r.WithContext(context.WithValue(r.Context(), "user_id", userId))

					next.ServeHTTP(w, r)
				})
			})

			protectedRouter.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status": "pong"}`))
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

			userssapi.HandlerFromMux(
				&users.Api{
					Repo: &userssqlite.Repo{
						Queries: *sqlc.New(db),
					},
					Txm: *infrastructure.NewTxManager(db),
				}, protectedRouter)

		})

	})

	//Frontend
	server.ServeFrontend(r)

	//Serve
	err = http.ListenAndServe(":8080", r)

	slog.Error("Error serving api", "error", err)

}
