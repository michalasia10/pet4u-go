package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"

	"src/internal/database"
	placepg "src/internal/modules/places/infrastructure/postgres"
)

func init() {
	goose.AddMigrationContext(upInitPlaces, downInitPlaces)
}

func upInitPlaces(ctx context.Context, _ *sql.Tx) error {
	m := database.Migrator()
	return m.AutoMigrate(&placepg.Place{})
}

func downInitPlaces(ctx context.Context, _ *sql.Tx) error {
	m := database.Migrator()
	return m.DropTable(&placepg.Place{})
}
