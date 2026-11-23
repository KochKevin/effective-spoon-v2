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

	
}
