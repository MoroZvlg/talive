package config

import (
	"os"
	"path"
	"strings"

	"github.com/joho/godotenv"
)

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	_ = godotenv.Load(path.Join(pwd, ".env"))
}

func LogLevel() string {
	return strings.ToLower(env("LOG_LEVEL", "INFO"))
}

func LogFormat() string {
	return strings.ToLower(env("LOG_FORMAT", "text"))
}

func env(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
