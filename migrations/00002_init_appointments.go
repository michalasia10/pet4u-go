package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"

	"src/internal/database"
	apptpg "src/internal/modules/appointments/infrastructure/postgres"
)

func init() {
	goose.AddMigrationContext(upInitAppointments, downInitAppointments)
}

func upInitAppointments(ctx context.Context, _ *sql.Tx) error {
	m := database.Migrator()
	return m.AutoMigrate(&apptpg.Appointment{})
}

func downInitAppointments(ctx context.Context, _ *sql.Tx) error {
	m := database.Migrator()
	return m.DropTable(&apptpg.Appointment{})
}
