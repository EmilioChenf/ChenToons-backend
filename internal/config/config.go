package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppPort    string
	AppEnv     string
	UploadDir  string
	RawDBURL   string
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
		AppPort:    leerEnvChen("PORT", leerEnvChen("APP_PORT", "8080")),
		AppEnv:     leerEnvChen("APP_ENV", "development"),
		UploadDir:  leerEnvChen("UPLOAD_DIR", "./uploads"),
		RawDBURL:   leerEnvChen("DATABASE_URL", ""),
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
	if c.RawDBURL != "" {
		return c.RawDBURL
	}
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
