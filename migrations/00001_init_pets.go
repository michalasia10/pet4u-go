package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"

	"src/internal/database"
	petpg "src/internal/modules/pets/infrastructure/postgres"
)

func init() {
	goose.AddMigrationContext(upInitPets, downInitPets)
}

func upInitPets(ctx context.Context, _ *sql.Tx) error {
	m := database.Migrator()
	if err := m.AutoMigrate(&petpg.Pet{}); err != nil {
		return err
	}
	if err := m.AutoMigrate(&petpg.MedicalRecord{}); err != nil {
		return err
	}
	return nil
}

func downInitPets(ctx context.Context, _ *sql.Tx) error {
	m := database.Migrator()
	if err := m.DropTable(&petpg.MedicalRecord{}); err != nil {
		return err
	}
	if err := m.DropTable(&petpg.Pet{}); err != nil {
		return err
	}
	return nil
}
