package config

import (
	"shared/utils"
)

type ServerConfig struct {
	Port string
}

func LoadServerConfig() ServerConfig {
	return ServerConfig{
		Port: utils.GetEnv("SERVER_PORT", "7070"),
	}
}
