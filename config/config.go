package config

import (
	"os"

	"github.com/sirjager/gopkg/db"
	"github.com/sirjager/gopkg/mail"
	"github.com/spf13/viper"

	"github.com/sirjager/trueauth/logger"
)

// Config represents the application configuration.
type Config struct {
	Database db.Config
	Mail     mail.Config
	Logger   logger.Config
	Server   ServerConfig
	Auth     AuthConfig
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

	if err = viper.Unmarshal(&config.Logger); err != nil {
		return
	}

	if err = viper.Unmarshal(&config.Database); err != nil {
		return
	}

	config.Database.Migrations = "file://" + config.Database.Migrations

	// Construct the DBUrl using the DBConfig values.
	if config.Server.ServerName == "" {
		config.Server.ServerName, _ = os.Hostname()
		config.Logger.ServerName, _ = os.Hostname()
	}

	return
}
