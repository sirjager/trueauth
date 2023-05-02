package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/sirjager/trueauth/cfg"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func Migrate(logger zerolog.Logger, conn *sql.DB, config cfg.DBConfig) (err error) {
	var driver database.Driver
	switch config.DBDriver {
	case "postgres", "postgresql":
		driver, err = postgres.WithInstance(conn, &postgres.Config{})
	case "mysql":
		driver, err = mysql.WithInstance(conn, &mysql.Config{})
	default:
		err = fmt.Errorf("'%s' is either not supported or not implemented", config.DBDriver)
	}
	if err != nil {
		return err
	}

	dbmigrate, err := migrate.NewWithDatabaseInstance(config.DBMigrate, config.DBDriver, driver)
	if err != nil {
		return err
	}

	err = dbmigrate.Up()
	if err != nil {
		if err != migrate.ErrNoChange {
			return err
		}
		logger.Info().Msg("database migration is up to date")
	} else {
		logger.Info().Msg("database migration complete")
	}

	return nil
}

func PingRedis(address string, logger zerolog.Logger) error {
	client := redis.NewClient(&redis.Options{Addr: address})
	defer client.Close()
	return client.Ping(context.Background()).Err()
}
