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
		config.AddConfigPath("../")    // For tests running from subdirectories
		config.AddConfigPath("../../") // For tests running from deeper subdirectories
		config.AutomaticEnv()
		config.WatchConfig()

		if err := config.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				// Config file not found - use defaults and environment variables
				// This is acceptable for tests or when using env vars only
				log.Printf("Config file not found, using defaults and environment variables: %v", err)
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
