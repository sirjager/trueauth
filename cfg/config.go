package cfg

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	StartTime   time.Time
	ServiceName string

	RestPort string `mapstructure:"REST_PORT"` //? port for serving "Rest" requests
	GrpcPort string `mapstructure:"GRPC_PORT"` //? port for serving "Grpc" requests

	TokenSecret         string        `mapstructure:"TOKEN_SECRET"`          //? access token time to live
	AccessTokenTTL      time.Duration `mapstructure:"ACCESS_TOKEN_TTL"`      //? access token time to live
	RefreshTokenTTL     time.Duration `mapstructure:"REFRESH_TOKEN_TTL"`     //? refres token time to live
	VerifyTokenTTL      time.Duration `mapstructure:"VERIFY_TOKEN_TTL"`      //? verification token time to live
	VerifyTokenCooldown time.Duration `mapstructure:"VERIFY_TOKEN_COOLDOWN"` //? verification token request cooldown
	ResetTokenTTL       time.Duration `mapstructure:"RESET_TOKEN_TTL"`       //? reset token time to live
	ResetTokenCooldown  time.Duration `mapstructure:"RESET_TOKEN_COOLDOWN"`  //? reset token request cooldown
	DeleteTokenTTL      time.Duration `mapstructure:"DELETE_TOKEN_TTL"`      //? reset token time to live
	DeleteTokenCooldown time.Duration `mapstructure:"DELETE_TOKEN_COOLDOWN"` //? reset token request cooldown

	DBConfig DBConfig //? database configs

	GmailSMTP GmailSMTP //? Gmail Stmp server configs
}

type DBConfig struct {
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
}

type GmailSMTP struct {
	SMTPSender string `mapstructure:"GMAIL_SMTP_NAME"` //? smtp account holder name
	SMTPUser   string `mapstructure:"GMAIL_SMTP_USER"` //? smtp email provider user
	SMTPPass   string `mapstructure:"GMAIL_SMTP_PASS"` //? smtp email provider pass
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

	if err = viper.Unmarshal(&config.DBConfig); err != nil {
		return
	}
	config.DBConfig.DBUrl = fmt.Sprintf("%s://%s:%s@%s:%s/%s%s", config.DBConfig.DBDriver, config.DBConfig.DBUser, config.DBConfig.DBPass, config.DBConfig.DBHost, config.DBConfig.DBPort, config.DBConfig.DBName, config.DBConfig.DBArgs)
	config.DBConfig.DBMigrate = "file://" + config.DBConfig.DBMigrate

	if err = viper.Unmarshal(&config.GmailSMTP); err != nil {
		return
	}

	return
}
