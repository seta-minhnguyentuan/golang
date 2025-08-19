package config

type ServerConfig struct {
	Port string
}

func LoadServerConfig() ServerConfig {
	return ServerConfig{
		Port: getEnv("SERVER_PORT", "8080"),
	}
}
