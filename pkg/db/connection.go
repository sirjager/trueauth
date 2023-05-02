package db

import (
	"database/sql"

	"github.com/rs/zerolog"
)

type Config struct {
	Name    string `mapstructure:"DB_NAME"`    // The name of the database
	Host    string `mapstructure:"DB_HOST"`    // The hostname or IP address of the database server
	Port    string `mapstructure:"DB_PORT"`    // The port number on which the database server is listening
	User    string `mapstructure:"DB_USER"`    // The username for authenticating with the database server
	Pass    string `mapstructure:"DB_PASS"`    // The password for authenticating with the database server
	Args    string `mapstructure:"DB_ARGS"`    // Additional arguments for the database connection
	Driver  string `mapstructure:"DB_DRIVER"`  // The database driver to use: postgres, mysql
	Migrate string `mapstructure:"DB_MIGRATE"` // The path to the database migration files
	Url     string // The URL for the database connection (derived from other fields)

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
	return &Database{conn, config, logr}, conn, err
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
