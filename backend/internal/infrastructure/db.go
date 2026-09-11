package infrastructure

import "embed"

//go:embed db/migrations
var DBMigrations embed.FS

//go:embed db/queries
var DBQueries embed.FS