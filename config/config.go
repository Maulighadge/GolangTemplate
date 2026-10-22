package config

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName            string
	AppEnv             string
	AppPort            string
	AppURL             string
	DBDriver           string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	DBSSLMode          string
	DBTimezone         string
	DBSQLitePath       string
	JWTSecret          string
	JWTExpirationHours int
	CORSAllowedOrigins string
}

var AppConfig *Config

// LoadConfig loads environment variables and parses into Config struct
func LoadConfig() *Config {
	// Try loading .env file if it exists (ignore error if missing, env vars will be used)
	if err := godotenv.Load(); err != nil {
		log.Println("[INFO] No .env file found or error loading, reading from system environment")
	}

	jwtExp, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "72"))
	if err != nil {
		jwtExp = 72
	}

	AppConfig = &Config{
		AppName:            getEnv("APP_NAME", "Golang-MVC-API"),
		AppEnv:             getEnv("APP_ENV", "development"),
		AppPort:            getEnv("APP_PORT", "8080"),
		AppURL:             getEnv("APP_URL", "http://localhost:8080"),
		DBDriver:           getEnv("DB_DRIVER", "sqlite"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", "postgres"),
		DBName:             getEnv("DB_NAME", "golang_template_db"),
		DBSSLMode:          getEnv("DB_SSLMODE", "disable"),
		DBTimezone:         getEnv("DB_TIMEZONE", "UTC"),
		DBSQLitePath:       getEnv("DB_SQLITE_PATH", "./app.db"),
		JWTSecret:          getEnv("JWT_SECRET", "default-jwt-secret-key-change-me"),
		JWTExpirationHours: jwtExp,
		CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "*"),
	}

	return AppConfig
}

// Validate validates that critical configuration values are safe and present
func (c *Config) Validate() error {
	if c.AppEnv == "production" {
		if c.JWTSecret == "default-jwt-secret-key-change-me" || len(c.JWTSecret) < 16 {
			return errors.New("JWT_SECRET must be explicitly set to a secure string of at least 16 characters in production")
		}
	}
	return nil
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}
