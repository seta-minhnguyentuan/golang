package config

type ServerConfig struct {
	Port string
}

func LoadServerConfig() ServerConfig {
	return ServerConfig{
		Port: getEnv("HTTP_PORT", "8080"),
	}
}
