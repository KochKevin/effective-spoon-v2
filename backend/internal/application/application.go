package application

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/KochKevin/effective-spoon-v2/internal/config"
	"github.com/KochKevin/effective-spoon-v2/internal/infrastructure"
	sqlc "github.com/KochKevin/effective-spoon-v2/internal/infrastructure/sqlite/generated"
)

type Application struct {
	Txm      *infrastructure.TxManager
	Repos    Repositories
	Services Services
	Route    http.Handler

	cfg config.Config
}

func New(ctx context.Context, db *sql.DB, cfg config.Config) (app *Application, err error) {

	//Wiring
	queries := sqlc.New(db)
	txm := infrastructure.NewTxManager(db)
	app = &Application{
		Txm:   txm,
		Repos: newRepos(queries),
		cfg:   cfg,
	}
	app.Services, err = app.newServices(ctx)
	if err != nil {
		return nil, fmt.Errorf("in creating app: %w", err)
	}

	//Routing
	app.Route = app.newRoutes()

	return app, nil
}

func (a *Application) Run() {

	err := http.ListenAndServe(":8080", a.Route)
	slog.Error("Error serving api", "error", err)
}
