package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

func InitializeViper() {
	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	viper.WatchConfig()

	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Print("Viper config change registered!")
	})
}
