package config

import "os"

type Config struct {
	ServerAddr   string
	DatabaseURL  string
	RedisAddr    string
	AdminToken   string
}

func Load() *Config {
	return &Config{
		ServerAddr:  getEnv("SERVER_ADDR", ":8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://vault:vault@localhost:5432/securevault?sslmode=disable"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		AdminToken:  getEnv("ADMIN_TOKEN", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
