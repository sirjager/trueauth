package db

import (
	"database/sql"

	"github.com/rs/zerolog"
)

type Config struct {
	Url       string `mapstructure:"DB_URL"`     // The database driver to use: postgres, mysql
	Driver    string `mapstructure:"DB_DRIVER"`  // Additional arguments for the database connection
	Migrate   string `mapstructure:"DB_MIGRATE"` // The path to the database migration files
	DBName    string `mapstructure:"DB_NAME"`    // The path to the database migration files
	User      string `mapstructure:"DB_USER"`    // The path to the database migration files
	Pass      string `mapstructure:"DB_PASS"`    // The path to the database migration files
	SSLMode   string `mapstructure:"DB_SSLMODE"` // The path to the database migration files
	RedisAddr string `mapstructure:"REDIS_ADDR"` // Redis connection string async workers
}

// Database represents a database connection.
type Database struct {
	conn   *sql.DB
	config Config
	logr   zerolog.Logger
}

// NewConnection creates a new database connection based on the provided configuration.
// It returns a Database instance and any error encountered during the connection process.
func NewDatabae(config Config, logr zerolog.Logger) (*Database, *sql.DB, error) {
	conn, err := sql.Open(config.Driver, config.Url)
	if err != nil {
		return nil, nil, err
	}
	// Return a new Database instance with the connection.
	return &Database{
		conn:   conn,
		logr:   logr,
		config: config,
	}, conn, err
}

// Close closes the database connection.
func (d *Database) Close() error {
	// Close the underlying database connection.
	return d.conn.Close()
}

// Ping database
func (d *Database) Ping() error {
	// Close the underlying database connection.
	return d.conn.Ping()
}
