package config

import (
	"log"
	"os"
)

// Config holds all configuration values loaded from environment variables
type Config struct {
	PocketbaseURL      string
	PocketbaseAdmin    string
	PocketbasePassword string
}

// Load reads and returns the configuration from environment variables
func Load() Config {
	cfg := Config{
		PocketbaseURL:      os.Getenv("PB_URL"),
		PocketbaseAdmin:    os.Getenv("PB_ADMIN_EMAIL"),
		PocketbasePassword: os.Getenv("PB_ADMIN_PASSWORD"),
	}

	// Fail fast if any required config is missing
	if cfg.PocketbaseURL == "" || cfg.PocketbaseAdmin == "" || cfg.PocketbasePassword == "" {
		log.Fatal("Missing one or more required environment variables: PB_URL, PB_ADMIN_EMAIL, PB_ADMIN_PASSWORD (should be teddy@scroobl.es for both email and password)")
	}

	return cfg
}
