package cfg

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	StartTime   time.Time
	ServiceName string

	RestPort string `mapstructure:"REST_PORT"` //? port for serving "Rest" requests
	GrpcPort string `mapstructure:"GRPC_PORT"` //? port for serving "Grpc" requests

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
	return
}
