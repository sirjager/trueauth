package config

import "time"

type ServerConfig struct {

	// StartTime is the timestamp when the application started.
	StartTime time.Time

	// AppName is the name of the service/app/product/company
	AppName string `mapstructure:"APP_NAME"`

	// ServiceName is the name of the service.
	ServerName string `mapstructure:"SERVER_NAME"`

	// Host is the hostname for the REST server.
	Host string `mapstructure:"HOST"`

	// Port is the port number for the REST server.
	Port string `mapstructure:"PORT"`

	// RestPort is the port number for the REST server.
	RestPort string `mapstructure:"REST_PORT"`

	// gRPC Port is the port number for the gRPC server.
	GrpcPort string `mapstructure:"GRPC_PORT"`

	MaxRequestTimeout time.Duration `mapstructure:"MAX_REQUEST_TIMEOUT"`
}
