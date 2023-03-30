package cfg

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	LogErrors bool `mapstructure:"LOG_ERRORS"` //? for logging errors in console

	RestPort string `mapstructure:"REST_PORT"` //? port for serving "Rest" requests
	GrpcPort string `mapstructure:"GRPC_PORT"` //? port for serving "Grpc" requests

	TokenSecret         string        `mapstructure:"TOKEN_SECRET"`          //? access token time to live
	AccessTokenTTL      time.Duration `mapstructure:"ACCESS_TOKEN_TTL"`      //? access token time to live
	RefreshTokenTTL     time.Duration `mapstructure:"REFRESH_TOKEN_TTL"`     //? refres token time to live
	VerifyTokenTTL      time.Duration `mapstructure:"VERIFY_TOKEN_TTL"`      //? verification token time to live
	VerifyTokenCooldown time.Duration `mapstructure:"VERIFY_TOKEN_COOLDOWN"` //? verification token request cooldown
	ResetTokenTTL       time.Duration `mapstructure:"RESET_TOKEN_TTL"`       //? reset token time to live
	ResetTokenCooldown  time.Duration `mapstructure:"RESET_TOKEN_COOLDOWN"`  //? reset token request cooldown

	RedisAddr string `mapstructure:"REDIS_ADDR"` //? Redis connection string for caching
	DBName    string `mapstructure:"DB_NAME"`    //? database name
	DBHost    string `mapstructure:"DB_HOST"`    //? database host
	DBPort    string `mapstructure:"DB_PORT"`    //? database port
	DBUser    string `mapstructure:"DB_USER"`    //? database user
	DBPass    string `mapstructure:"DB_PASS"`    //? database pass
	DBArgs    string `mapstructure:"DB_ARGS"`    //? database args
	DBDriver  string `mapstructure:"DB_DRIVER"`  //? database driver: postgres, mysql
	DBMigrate string `mapstructure:"DB_MIGRATE"` //? database migation files path
	DBUrl     string //? database url

	TestDBUrl string `mapstructure:"TEST_DB_URL"` //? test database url

	SMTPHost  string `mapstructure:"SMTP_HOST"`  //? smtp email provider host
	SMTPPort  string `mapstructure:"SMTP_PORT"`  //? smtp email provider port
	SMTPUser  string `mapstructure:"SMTP_USER"`  //? smtp email provider user
	SMTPPass  string `mapstructure:"SMTP_PASS"`  //? smtp email provider pass
	SMTPEmail string `mapstructure:"SMTP_EMAIL"` //? will use this email to send all emails

	StartTime   time.Time
	ServiceName string
}

func LoadConfigs(path, name string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	config.DBUrl = fmt.Sprintf("%s://%s:%s@%s:%s/%s%s", config.DBDriver, config.DBUser, config.DBPass, config.DBHost, config.DBPort, config.DBName, config.DBArgs)
	config.DBMigrate = "file://" + config.DBMigrate
	return
}
