package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type Config struct {
	Migrations string `mapstructure:"MIGRATIONS_DIR"`  // The path to the database migration files
	Driver     string `mapstructure:"DATABASE_DRIVER"` // The database driver to use: postgres, mysql
	Name       string `mapstructure:"DATABASE_NAME"`   // The name of the database
	Host       string `mapstructure:"DATABASE_HOST"`   // The hostname or IP address of the database server
	Port       string `mapstructure:"DATABASE_PORT"`   // The port number on which the database server is listening
	User       string `mapstructure:"DATABASE_USER"`   // The username for authenticating with the database server
	Pass       string `mapstructure:"DATABASE_PASS"`   // The password for authenticating with the database server
	Args       string `mapstructure:"DATABASE_ARGS"`   // Additional arguments for the database connection
	URL        string // The URL for the database connection (derived from other fields)
	RedisURL   string `mapstructure:"REDIS_URL"`      // The URL for the redis connction
	RedisAddr  string `mapstructure:"REDIS_ADDRESS"`  // The URL for the redis connction
	RedisUser  string `mapstructure:"REDIS_USERNAME"` // The username for authenticating with the redis server
	RedisPass  string `mapstructure:"REDIS_PASSWORD"` // The password for authenticating with the redis server
	RedisPort  string `mapstructure:"REDIS_PORT"`     // The port number on which the redis server is listening
}

// Database represents a database connection.
type Database struct {
	logr   zerolog.Logger
	pool   *pgxpool.Pool
	config Config
}

// NewConnection creates a new database connection based on the provided configuration.
// It returns a Database instance and any error encountered during the connection process.
func NewDatabae(
	ctx context.Context,
	config Config,
	logr zerolog.Logger,
) (*Database, *pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, config.URL)
	if err != nil {
		return nil, nil, err
	}
	// Return a new Database instance with the connection.
	return &Database{logr, pool, config}, pool, err
}

// Close closes the database connection.
func (database *Database) Close() {
	database.pool.Close()
}
