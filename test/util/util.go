package util

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
}

func LoadTestConfig() (*Config, error) {
	os.Setenv("APP_ENVIRONMENT", "test")

	config := &Config{
		DBHost:     getEnvOrDefault("DATABASE_HOST", "localhost"),
		DBPort:     getEnvOrDefault("DATABASE_PORT", "5432"),
		DBUser:     getEnvOrDefault("POSTGRES_USER", "root"),
		DBPassword: getEnvOrDefault("POSTGRES_PASSWORD", "lets-jungle-it-bro!"),
		DBName:     "go-auth-db-test",
		JWTSecret:  "test-secret",
	}

	viper.SetConfigName("config.test")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../configs")

	viper.Set("jwt.secret", config.JWTSecret)
	viper.Set("jwt.access_token_expiry", "15m")
	viper.Set("jwt.refresh_token_expiry", "24h")

	return config, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetTestDSN(config *Config) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBName,
		config.DBPassword,
	)
}
