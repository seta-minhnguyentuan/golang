package utils

import (
	"log"
	"os"
	"strconv"
	"time"
)

func MustEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		if def != "" {
			return def
		}
		log.Fatalf("missing required env %s", key)
	}
	return v
}

func AsInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("bad int for %s: %v", key, err)
	}
	return i
}

func AsInt64(key string, def int64) int64 {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		log.Fatalf("bad int64 for %s: %v", key, err)
	}
	return i
}

func AsDuration(key string, def time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		log.Fatalf("bad duration for %s: %v", key, err)
	}
	return d
}

func GetEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
