package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

var (
	config *viper.Viper
	once   sync.Once
)

// GetConfig returns the singleton config instance, initializing it on first call
func GetConfig() *viper.Viper {
	once.Do(func() {
		config = viper.New()
		config.SetConfigName("config")
		config.SetConfigType("yaml")
		config.AddConfigPath("./")
		config.AutomaticEnv()
		config.WatchConfig()

		if err := config.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				log.Fatalf("Config file not found: %v", err)
			} else {
				log.Fatalf("Error reading config file: %v", err)
			}
		}
	})
	return config
}

// GetString returns a string value from the config
func GetString(key string) string {
	return GetConfig().GetString(key)
}
