package postgres

import (
	"fmt"
	"os"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func LoadConfigFromEnv() Config {
	return Config{
		Host:     getEnv("PGHOST", "localhost"),
		Port:     getEnv("PGPORT", "5432"),
		User:     getEnv("PGUSER", "postgres"),
		Password: os.Getenv("PGPASSWORD"),
		Database: getEnv("PGDATABASE", "postgres"),
		SSLMode:  getEnv("PGSSLMODE", "disable"),
	}
}

func (c Config) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.SSLMode,
	)
}

func (c Config) IsConfigured() bool {
	return os.Getenv("DATABASE_URL") != "" || os.Getenv("PGHOST") != ""
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
