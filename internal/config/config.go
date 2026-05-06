package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppPort    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	CORSOrigin string
}

func LoadCheninConfig() Config {
	return Config{
		AppPort:    leerEnvChen("APP_PORT", "8080"),
		DBHost:     leerEnvChen("DB_HOST", "localhost"),
		DBPort:     leerEnvChen("DB_PORT", "5432"),
		DBUser:     leerEnvChen("DB_USER", "chentoons"),
		DBPassword: leerEnvChen("DB_PASSWORD", "chentoons123"),
		DBName:     leerEnvChen("DB_NAME", "chentoons_db"),
		DBSSLMode:  leerEnvChen("DB_SSLMODE", "disable"),
		CORSOrigin: leerEnvChen("CORS_ORIGIN", "*"),
	}
}

func (c Config) DatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
		c.DBSSLMode,
	)
}

func leerEnvChen(nombre string, fallback string) string {
	valor := os.Getenv(nombre)
	if valor == "" {
		return fallback
	}
	return valor
}
