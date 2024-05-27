package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	"github.com/sirjager/trueauth/pkg/db"
	"github.com/sirjager/trueauth/pkg/mail"
)

// Config represents the application configuration.
type Config struct {
	Server   ServerConfig // holds the configuration for server itself
	Database db.Config    // holds the configuration for the database.
	Mail     mail.Config  // holds the configuration for the email smtp server.
	Auth     AuthConfig   // holds the configuration for the auth.
}

// LoadConfigs loads the configuration from the specified YAML file.
func LoadConfigs(path string, name string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	if err = viper.Unmarshal(&config.Server); err != nil {
		return
	}

	if err = viper.Unmarshal(&config.Auth); err != nil {
		return
	}

	if err = viper.Unmarshal(&config.Mail); err != nil {
		return
	}

	if err = viper.Unmarshal(&config.Database); err != nil {
		return
	}

	config.Database.Migrations = "file://" + config.Database.Migrations

	// Construct the DBUrl using the DBConfig values.
	config.Database.URL = fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s%s",
		config.Database.Driver,
		config.Database.User,
		config.Database.Pass,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
		config.Database.Args,
	)

	if config.Server.ServerName == "" {
		config.Server.ServerName, _ = os.Hostname()
	}

	return
}
