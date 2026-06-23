package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {

	r := chi.NewRouter()
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

	http.ListenAndServe(":5000", r)

}
