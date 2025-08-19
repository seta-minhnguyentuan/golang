package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	DBHost    string
	DBUser    string
	DBPass    string
	DBName    string
	DBPort    int
	DBSSLMode string
}

var cfg *DatabaseConfig

func LoadDB() *DatabaseConfig {
	_ = godotenv.Load()

	if cfg != nil {
		return cfg
	}

	port, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		log.Fatalf("Invalid DB_PORT: %v", err)
	}

	cfg = &DatabaseConfig{
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBUser:    getEnv("DB_USER", "postgres"),
		DBPass:    getEnv("DB_PASSWORD", "postgres"),
		DBName:    getEnv("DB_NAME", "user_service"),
		DBPort:    port,
		DBSSLMode: getEnv("DB_SSLMODE", "disable"),
	}

	return cfg
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		c.DBHost, c.DBUser, c.DBPass, c.DBName, c.DBPort, c.DBSSLMode,
	)
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
