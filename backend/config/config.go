package config

import (
	"os"
)

type Config struct {
	Port       string
	DBFilePath string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./points.db"
	}

	return &Config{
		Port:       port,
		DBFilePath: dbPath,
	}
}
