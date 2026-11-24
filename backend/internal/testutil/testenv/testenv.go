package testenv

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"net/http/httptest"
	"testing"

	"github.com/KochKevin/effective-spoon-v2/internal/application"
	"github.com/KochKevin/effective-spoon-v2/internal/config"
	"github.com/KochKevin/effective-spoon-v2/internal/infrastructure"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

type Env struct {
	DB  *sql.DB
	URL string
	App *application.Application
}

func New(t *testing.T) *Env {

	t.Helper()

	logger := slog.New(slog.NewJSONHandler(t.Output(), &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	db := openDB(t)
	migrate(t, db, logger)

	app, err := application.New(t.Context(), db, config.Config{StripeClientKey: "NO_KEY"})
	if err != nil {
		slog.Error("Error in app", "error", fmt.Errorf("in building app: %w", err))
		return nil
	}

	srv := httptest.NewServer(app.Route)
	t.Cleanup(srv.Close)

	return &Env{DB: db, URL: srv.URL, App: app}

}

func openDB(t *testing.T) *sql.DB {
	t.Helper()

	sqliteConnectionString := "file::memory:?_pragma=journal_mode=WAL&_pragma=busy_timeout=5000"

	db, err := goose.OpenDBWithDriver("sqlite", sqliteConnectionString)
	require.NoError(t, err, "no error opening db")
	if err != nil {
		slog.Error("can not open database", "err", err)
		return nil
	}
	t.Cleanup(func() { _ = db.Close() })

	return db
}

func migrate(t *testing.T, db *sql.DB, logger *slog.Logger) {
	t.Helper()

	migrationsFS, err := fs.Sub(infrastructure.DBMigrations, "db/migrations")
	require.NoError(t, err, "no error opening migration sub filesystem")
	if err != nil {
		slog.Error("failed to create sub filesystem for migrations", "err", err)
		return
	}

	gooseProvider, err := goose.NewProvider(goose.DialectSQLite3, db, migrationsFS, goose.WithSlog(logger))
	require.NoError(t, err, "no error creating goose provider")
	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = gooseProvider.Up(t.Context())
	require.NoError(t, err, "no error doing goose migration")
	if err != nil {
		slog.Error("Error migrating database", "error", err)
		return
	}

}
