package configs

import "os"

type Config struct {
	AppPort   string
	MongoUri  string
	MongoDB   string
	JWTSecret string
}

func LoadConfig() (*Config, error) {

	var cfg Config

	if val := os.Getenv("APP_PORT"); val != "" {
		cfg.AppPort = val
	}

	if val := os.Getenv("MONGO_URI"); val != "" {
		cfg.MongoUri = val
	}

	if val := os.Getenv("MONGO_DB"); val != "" {
		cfg.MongoDB = val
	}

	if val := os.Getenv("JWT_SECRET"); val != "" {
		cfg.JWTSecret = val
	}

	return &cfg, nil
}