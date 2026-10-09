package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort       string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBSSLMode     string
	JWTSecret     string
	AdminName     string
	AdminEmail    string
	AdminPassword string
}

func LoadConfig() (*Config, error) {
	// Load .env file
	_ = godotenv.Load()

	cfg := &Config{
		AppPort:       os.Getenv("APP_PORT"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		DBSSLMode:     os.Getenv("DB_SSLMODE"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		AdminName:     os.Getenv("ADMIN_NAME"),
		AdminEmail:    os.Getenv("ADMIN_EMAIL"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
	}

	if cfg.AppPort == "" {
		cfg.AppPort = "8080"
	}

	if cfg.DBSSLMode == "" {
		cfg.DBSSLMode = "disable"
	}

	if cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBUser == "" || cfg.DBName == "" {
		return nil, fmt.Errorf("database configuration is incomplete")
	}

	if len(cfg.JWTSecret) < 32 {
		return nil, fmt.Errorf("JWT_SECRET must be at least 32 characters")
	}

	if cfg.AdminName == "" ||
		cfg.AdminEmail == "" ||
		cfg.AdminPassword == "" {
		return nil, fmt.Errorf("admin seed configuration missing")
	}
	return cfg, nil
}
