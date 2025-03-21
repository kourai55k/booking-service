package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env                string `yaml:"env" env-default:"local"`
	PostgresConnString string `yaml:"PostgresConnString" env-required:"true"`
}

// MustLoad loads the configuration
func MustLoad() *Config {
	// Only load .env file if CONFIG_PATH is not set (assumes running locally)
	if os.Getenv("CONFIG_PATH") == "" {
		if err := godotenv.Load(); err != nil {
			log.Println("Warning: No .env file found, using system environment variables")
		}
	}

	// Get config path from environment
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	var cfg Config

	// Read config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Cannot read config: %s", err)
	}

	return &cfg
}
