package migrations

import (
	"database/sql"
)

// This is just an example how to define Go migrations.
func init() {
	goose.AddMigration(upNoop, downNoop)
}

func upNoop(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return nil
}

func downNoop(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil // migrate.ErrDownNotSupported
}
