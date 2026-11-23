package application

import (
	"log/slog"
	"net/http"

	"github.com/KochKevin/effective-spoon-v2/internal/auth"
	authapi "github.com/KochKevin/effective-spoon-v2/internal/auth/generated"
	authservice "github.com/KochKevin/effective-spoon-v2/internal/auth/service"
	"github.com/KochKevin/effective-spoon-v2/internal/chargements"
	chargementsapi "github.com/KochKevin/effective-spoon-v2/internal/chargements/generated"
	"github.com/KochKevin/effective-spoon-v2/internal/events"
	eventsapi "github.com/KochKevin/effective-spoon-v2/internal/events/generated"
	"github.com/KochKevin/effective-spoon-v2/internal/input"
	inputapi "github.com/KochKevin/effective-spoon-v2/internal/input/generated"
	"github.com/KochKevin/effective-spoon-v2/internal/products"
	productsapi "github.com/KochKevin/effective-spoon-v2/internal/products/generated"
	"github.com/KochKevin/effective-spoon-v2/internal/push"
	pushapi "github.com/KochKevin/effective-spoon-v2/internal/push/generated"
	"github.com/KochKevin/effective-spoon-v2/internal/server"
	"github.com/KochKevin/effective-spoon-v2/internal/shoppingcarts"
	shoppingcartssapi "github.com/KochKevin/effective-spoon-v2/internal/shoppingcarts/generated"
	"github.com/KochKevin/effective-spoon-v2/internal/users"
	userssapi "github.com/KochKevin/effective-spoon-v2/internal/users/generated"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (a *Application) newRoutes() chi.Router {

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

	r.Route("/api", a.apiRoute)

	//Frontend
	server.ServeFrontend(r)

	return r

}

func (a *Application) apiRoute(r chi.Router) {

	authapi.HandlerFromMux(
		&auth.Api{
			UserRepo:    a.Repos.Users,
			AuthService: a.Services.Auth,
			Txm:         a.Txm,
		}, r)

	inputapi.HandlerFromMux(
		&input.Api{
			InputService: a.Services.Input,
		}, r)

	//Server Sent Events, Push Api endpoint
	pushapi.HandlerFromMux(
		&push.Api{
			PushService: a.Services.Push,
		}, r)

	r.Group(a.protectedApiRoute)
}

func (a *Application) protectedApiRoute(r chi.Router) {
	//GetLoggedInUser Middlewear
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Create/Get test user id
			//userId := uuid.Nil

			userId, err := a.Services.Auth.GetCurrentUserId()
			if err != nil {
				slog.Error("error: a user needs to be logged in to access this endpoint", "error", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			slog.Debug("auth middlewear", "user_id", userId)

			r = r.WithContext(authservice.WithUser(r.Context(), userId))

			next.ServeHTTP(w, r)
		})
	})

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status": "pong"}`))
	})

	productsapi.HandlerFromMux(
		&products.Api{
			Repo: a.Repos.Products,
			Txm:  a.Txm,
		}, r)

	shoppingcartssapi.HandlerFromMux(
		&shoppingcarts.Api{
			Service: a.Services.Shoppingcart,
		}, r)

	//Users
	userssapi.HandlerFromMux(
		&users.Api{
			Service: a.Services.User,
		}, r)

	chargementsapi.HandlerFromMux(&chargements.Api{
		Service: a.Services.Chargement,
	}, r)

	//Events handler
	eventsapi.HandlerFromMux(&events.Api{
		Service:     a.Services.Event,
		UserService: a.Services.User,
		Usage:       a.Services.EventUsage,
		Lifecycle:   a.Services.EventLifecycle,
	}, r)

}
