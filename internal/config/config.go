package config

import (
	"os"
)

type Config struct {
	DatabaseURL string
	RedisURL    string
	AdminToken  string
	ServerPort  string
	JWTSecret   string
}

func Load() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://vault:vault_secret@localhost:5432/samsung_vault?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379/0"),
		AdminToken:  getEnv("ADMIN_TOKEN", ""),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		JWTSecret:   getEnv("JWT_SECRET", ""),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
