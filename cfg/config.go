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

	DBConfig DBConfig //? database configs
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

	var dbConfig DBConfig
	if err = viper.Unmarshal(&dbConfig); err != nil {
		return
	}
	dbConfig.DBUrl = fmt.Sprintf("%s://%s:%s@%s:%s/%s%s", dbConfig.DBDriver, dbConfig.DBUser, dbConfig.DBPass, dbConfig.DBHost, dbConfig.DBPort, dbConfig.DBName, dbConfig.DBArgs)
	dbConfig.DBMigrate = "file://" + dbConfig.DBMigrate

	config.DBConfig = dbConfig

	return
}
