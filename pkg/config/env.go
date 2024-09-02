package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBSource         string
	ServerAddr       string
	CookieSecret     string
	CookieAge        time.Duration
	CookieIsSecure   bool
	CookieIsHttpOnly bool
}

var Envs = initConfig()

var fifteenMinutesInSeconds = time.Minute * 15

func initConfig() Config {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	return Config{
		DBSource:         getEnv("DB_SOURCE", ""),
		ServerAddr:       getEnv("SERVER_ADDR", "127.0.0.1:8080"),
		CookieSecret:     getEnv("COOKIE_SECRET", "A3256DCEEF11B26457FB1E8779552"),
		CookieAge:        getEnvAsDuration("COOKIE_AGE", fifteenMinutesInSeconds),
		CookieIsSecure:   getEnvAsBool("COOCKIE_IS_SECURE", true),
		CookieIsHttpOnly: getEnvAsBool("COOKIE_IS_HTTP_ONLY", true),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	value, found := os.LookupEnv(key)

	if !found {
		return fallback
	}

	durationValue, err := time.ParseDuration(value)

	if err != nil {
		return fallback
	}
	return durationValue
}

func getEnvAsInt(key string, fallback int) int {
	value, found := os.LookupEnv(key)

	if !found {
		return fallback
	}

	intValue, err := strconv.Atoi(value)

	if err != nil {
		return fallback
	}
	return intValue
}

func getEnvAsBool(key string, fallback bool) bool {
	value, found := os.LookupEnv(key)

	if !found {
		return fallback
	}

	boolValue, err := strconv.ParseBool(value)

	if err != nil {
		return fallback
	}
	return boolValue
}
