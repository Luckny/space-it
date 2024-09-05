package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variable.
type Config struct {
	DBSource         string        `mapstructure:"DB_SOURCE"`
	ServerAddr       string        `mapstructure:"SERVER_ADDR"`
	CookieSecret     string        `mapstructure:"COOKIE_SECRET"`
	CookieAge        time.Duration `mapstructure:"COOKIE_AGE"`
	CookieIsSecure   bool          `mapstructure:"COOKIE_IS_SECURE"`
	CookieIsHttpOnly bool          `mapstructure:"COOKIE_IS_HTTP_ONLY"`
}

// LoadConfig reads configuration from file or environment variables.
func Load(path string) (config Config) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		panic(err)
	}

	return
}
