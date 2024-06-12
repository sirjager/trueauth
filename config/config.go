package config

import (
	"os"
	"time"

	"github.com/sirjager/gopkg/db"
	"github.com/sirjager/gopkg/mail"
	"github.com/spf13/viper"

	"github.com/sirjager/trueauth/consul"
	"github.com/sirjager/trueauth/logger"
)

// Config represents the application configuration.
type Config struct {
	// StartTime is the timestamp when the application started.
	StartTime time.Time

	// ServiceName is the name of the service
	ServiceName string `mapstructure:"SERVICE_NAME"`

	// ServerName is the name of the server, could be the hostname, or container's hostname.
	ServerName string `mapstructure:"SERVER_NAME"`

	// Host is the hostname for the REST server.
	Host string `mapstructure:"HOST"`

	// Port is the port number for the REST server.
	Port int `mapstructure:"PORT"`

	// RestPort is the port number for the REST server.
	RestPort int `mapstructure:"REST_PORT"`

	// gRPC Port is the port number for the gRPC server.
	GrpcPort int `mapstructure:"GRPC_PORT"`

	MaxRequestTimeout time.Duration `mapstructure:"MAX_REQUEST_TIMEOUT"`

	Database db.Config
	Mail     mail.Config
	Logger   logger.Config
	Auth     AuthConfig
	Consul   consul.Config
}

// LoadConfigs loads the configuration from the specified YAML file.
func LoadConfigs(path string, name string, startTime time.Time) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	if err = viper.Unmarshal(&config); err != nil {
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

	if err = viper.Unmarshal(&config.Consul); err != nil {
		return
	}

	// Construct the DBUrl using the DBConfig values.
	if config.ServerName == "" {
		config.ServerName, _ = os.Hostname()
		config.Logger.ServerName = config.ServerName
		config.Consul.ServerName = config.ServerName
	}

	config.StartTime = startTime
	config.Consul.StartTime = startTime

	return
}
