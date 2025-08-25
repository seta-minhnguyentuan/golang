package config

import (
	"fmt"
	"log"
	"shared/utils"
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

	port, err := strconv.Atoi(utils.GetEnv("DB_PORT", "5432"))
	if err != nil {
		log.Fatalf("Invalid DB_PORT: %v", err)
	}

	cfg = &DatabaseConfig{
		DBHost:    utils.GetEnv("DB_HOST", "localhost"),
		DBUser:    utils.GetEnv("DB_USER", "postgres"),
		DBPass:    utils.GetEnv("DB_PASSWORD", ""),
		DBName:    utils.GetEnv("DB_NAME", "postgres"),
		DBPort:    port,
		DBSSLMode: utils.GetEnv("DB_SSLMODE", "disable"),
	}

	return cfg
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		c.DBHost, c.DBUser, c.DBPass, c.DBName, c.DBPort, c.DBSSLMode,
	)
}
