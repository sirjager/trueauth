package config

import (
	"fmt"
	"time"

	"github.com/sirjager/trueauth/pkg/db"
	"github.com/sirjager/trueauth/pkg/mail"
	"github.com/spf13/viper"
)

type Config struct {
	LogErrors bool `mapstructure:"LOG_ERRORS"` //? for logging errors in console

	//? for internal
	StartTime   time.Time // StartTime is the timestamp when the application started.
	ServiceName string    // ServiceName is the name of the service.

	GrpcPort            string        `mapstructure:"GRPC_PORT"`             // GrpcPort is the port number for the gRPC server.
	GatewayPort         string        `mapstructure:"GATEWAY_PORT"`          // RestPort is the port number for the REST server.
	TokenSecret         string        `mapstructure:"TOKEN_SECRET"`          // access token time to live
	AccessTokenTTL      time.Duration `mapstructure:"ACCESS_TOKEN_TTL"`      // access token time to live
	RefreshTokenTTL     time.Duration `mapstructure:"REFRESH_TOKEN_TTL"`     // refres token time to live
	VerifyTokenTTL      time.Duration `mapstructure:"VERIFY_TOKEN_TTL"`      // verification token time to live
	VerifyTokenCooldown time.Duration `mapstructure:"VERIFY_TOKEN_COOLDOWN"` // verification token request cooldown
	ResetTokenTTL       time.Duration `mapstructure:"RESET_TOKEN_TTL"`       // reset token time to live
	ResetTokenCooldown  time.Duration `mapstructure:"RESET_TOKEN_COOLDOWN"`  // reset token request cooldown
	DeleteTokenTTL      time.Duration `mapstructure:"DELETE_TOKEN_TTL"`      // reset token time to live
	DeleteTokenCooldown time.Duration `mapstructure:"DELETE_TOKEN_COOLDOWN"` // reset token request cooldown

	RedisAddr string `mapstructure:"REDIS_ADDR"` // Redis connection string async workers

	//? for pkg
	Mail     mail.Config // Mail holds the configuration for mail-related settings.
	Database db.Config   // DBConfig holds the configuration for the database.

}

func LoadConfigs(path, name string) (config Config, err error) {
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

	if err = viper.Unmarshal(&config.Mail); err != nil {
		return
	}

	if err = viper.Unmarshal(&config.Database); err != nil {
		return
	}

	// Construct the DBUrl using the DBConfig values.
	config.Database.Url = fmt.Sprintf("%s://%s:%s@%s:%s/%s%s", config.Database.Driver, config.Database.User, config.Database.Pass, config.Database.Host, config.Database.Port, config.Database.Name, config.Database.Args)
	config.Database.Migrate = "file://" + config.Database.Migrate
	return
}
