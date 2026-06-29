package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/KochKevin/effective-spoon-v2/internal/infrastructure"
	sqlc "github.com/KochKevin/effective-spoon-v2/internal/infrastructure/sqlite/generated"
	"github.com/KochKevin/effective-spoon-v2/internal/products"
	productsapi "github.com/KochKevin/effective-spoon-v2/internal/products/generated"
	productssqlite "github.com/KochKevin/effective-spoon-v2/internal/products/sqlite"
	"github.com/KochKevin/effective-spoon-v2/internal/shoppingcarts"
	shoppingcartssapi "github.com/KochKevin/effective-spoon-v2/internal/shoppingcarts/generated"
	shoppingcartssqlite "github.com/KochKevin/effective-spoon-v2/internal/shoppingcarts/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/pressly/goose/v3"

	_ "modernc.org/sqlite"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	slog.Info("Backend Api started")

	//Do Database migrations
	db, err := goose.OpenDBWithDriver("sqlite", "./db/data/data.db")

	if err != nil {
		log.Fatal(err)
		return
	}

	defer db.Close()

	goose.Up(db, "./db/migrations")
	

	//Router Setup
	r := chi.NewRouter()

	//Middlewear
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://localhost:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "pong"}`))
	})

	//Routes

	productsapi.HandlerFromMux(&products.Api{
		Repo: &productssqlite.Repo{
			Queries: *sqlc.New(db),
		},
		Txm: *infrastructure.NewTxManager(db),
	}, r)

	shoppingcartssapi.HandlerFromMux(&shoppingcarts.Api{
		Repo: &shoppingcartssqlite.Repo{
			Queries: *sqlc.New(db),
		},
		ProductRepo: &productssqlite.Repo{
			Queries: *sqlc.New(db),
		},
		Txm: *infrastructure.NewTxManager(db),
	}, r)

	//Serve
	http.ListenAndServe(":8080", r)

}
