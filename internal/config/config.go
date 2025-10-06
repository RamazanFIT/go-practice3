package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	APIKey          string
	LogLevel        string
	DefaultLanguage string
	ServerHost      string
	ServerPort      string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	cfg := &Config{
		APIKey:          getEnv("API_KEY", ""),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		DefaultLanguage: getEnv("DEFAULT_LANGUAGE", "en"),
		ServerHost:      getEnv("SERVER_HOST", "localhost"),
		ServerPort:      getEnv("SERVER_PORT", "8080"),
	}

	if cfg.APIKey == "" {
		log.Fatal("API_KEY is not set in environment")
	}

	return cfg
}

func (c *Config) Address() string {
	return c.ServerHost + ":" + c.ServerPort
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
