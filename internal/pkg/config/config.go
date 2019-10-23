package config

import (
	"github.com/op/go-logging"
	"github.com/spf13/viper"
)

var log = logging.MustGetLogger("config")

func init() {
	viper.SetConfigType("json")
	viper.SetConfigName("server")
	viper.AddConfigPath("configs")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	log.Info("Config file used: ", viper.ConfigFileUsed())
}
