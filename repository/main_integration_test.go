package repository

import (
	"native-free-pollings/config"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../.env"); err != nil {
		_ = godotenv.Load("../../.env")
	}

	code := m.Run()
	os.Exit(code)
}

func Get() *config.Config {
	return &config.Config{
		Server: config.Server{
			Host: os.Getenv("APP_HOST"),
			Port: os.Getenv("APP_PORT"),
		},
		Database: config.Database{
			Host: os.Getenv("DB_HOST"),
			Port: os.Getenv("DB_PORT"),
			Name: os.Getenv("DB_NAME"),
			User: os.Getenv("DB_USER"),
			Pass: os.Getenv("DB_PASS"),
			SSL:  os.Getenv("DB_SSLMODE"),
		},
		JwtKey: []byte(os.Getenv("JWT_KEY")),
	}
}
