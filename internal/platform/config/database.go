package config

import "os"

type DatabaseConfig struct {
	URL string
}

func LoadDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		URL: os.Getenv("DATABASE_URL"),
	}
}
