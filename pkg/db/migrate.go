package db

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// Migrate performs database migrations using the provided configuration.
// It takes a logger, a database connection, and a DBConfig object.
// It returns an error if any error occurs during migration.
func (d *Database) Migrate() (err error) {
	var driver database.Driver

	// Determine the database driver based on the DBConfig.DBDriver value.
	// Create the corresponding migrate database driver instance.
	switch d.config.Driver {
	case "postgres", "postgresql":
		driver, err = postgres.WithInstance(d.conn, &postgres.Config{})
	case "mysql":
		driver, err = mysql.WithInstance(d.conn, &mysql.Config{})
	default:
		err = fmt.Errorf("'%s' is either not supported or not implemented", d.config.Driver)
	}
	if err != nil {
		return err
	}

	// Create a new migrate instance with the provided database driver.
	dbmigrate, err := migrate.NewWithDatabaseInstance(d.config.Migrate, d.config.Driver, driver)
	if err != nil {
		return err
	}

	// Apply any pending migrations to the database.
	err = dbmigrate.Up()
	if err != nil {
		// Check if the error is "no change" which indicates that there are no pending migrations.
		// Log an info message in this case.
		if err != migrate.ErrNoChange {
			return err
		}
		d.logr.Info().Msg("database migration is up to date")
	} else {
		d.logr.Info().Msg("database migration complete")
	}
	return nil
}
