package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/KochKevin/effective-spoon-v2/internal/products"
	productsapi "github.com/KochKevin/effective-spoon-v2/internal/products/generated"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/pressly/goose/v3"

	_ "modernc.org/sqlite"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	slog.SetDefault(logger)

	slog.Info("Backend Api started")

	//Do Database migrations
	db, err := goose.OpenDBWithDriver("sqlite", "app/data/main_data.db")

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
		// WICHTIG FÜR CODESPACES:
		// Da sich deine Codespace-Frontend-URLs ständig ändern können,
		// erlaubt "AllowedOrigins" mit "*" oder Wildcards den Zugriff von überall während der Entwicklung.
		AllowedOrigins:   []string{"https://*", "http://localhost:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Cached die CORS-Antwort für 5 Minuten
	}))

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "pong"}`))
	})

	//Routes

	productsapi.HandlerFromMux(&products.Api{}, r)

	//Serve
	http.ListenAndServe(":8080", r)

}
