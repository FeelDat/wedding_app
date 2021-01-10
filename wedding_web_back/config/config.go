package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Configuration struct {
	Server ServerConfiguration
	AccessSecret string `json:"access_secret"`
	RefreshSecret string `json:"refresh_secret"`
}

type ServerConfiguration struct {
	Address string
	Port    int
}

func GetConfig(log *zap.Logger) *Configuration {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	viper.SetDefault("server.address", "127.0.0.1")
	viper.SetDefault("server.port", 18000)

	viper.SetConfigType("yml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Error(err.Error())
	}

	configuration := &Configuration{}
	err := viper.Unmarshal(&configuration)
	if err != nil {
		log.Error(err.Error())
	} else {
		log.Debug("Using config file: " + viper.ConfigFileUsed())
	}

	return configuration
}
