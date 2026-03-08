package env

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port         int64
		Host         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		IdleTimeout  time.Duration
	}

	Database struct {
		Host     string
		Port     int64
		User     string
		Password string
		DBName   string
		SSLMode  string
	}

	JWT struct {
		Secret        string
		TokenExpiry   time.Duration
		RefreshExpiry time.Duration
	}

	Environment string
}

func Load() (*Config, error) {
	godotenv.Load()

	cfg := &Config{}

	// Server config
	cfg.Server.Port = getEnvInt("SERVER_PORT", 4000)
	cfg.Server.Host = getEnv("SERVER_HOST", "0.0.0.0")
	cfg.Server.ReadTimeout = time.Second * 15
	cfg.Server.WriteTimeout = time.Second * 15
	cfg.Server.IdleTimeout = time.Second * 60

	// Database config
	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnvInt("DB_PORT", 5432)
	cfg.Database.User = getEnv("DB_USER", "postgres")
	cfg.Database.Password = getEnv("DB_PASSWORD", "")
	cfg.Database.DBName = getEnv("DB_NAME", "auth_service")
	cfg.Database.SSLMode = getEnv("DB_SSLMODE", "disable")

	// JWT config
	cfg.JWT.Secret = getEnv("JWT_SECRET", "default-secret-key")
	cfg.JWT.TokenExpiry = time.Hour * 24    // 24 hours
	cfg.JWT.RefreshExpiry = time.Hour * 168 // 7 days

	cfg.Environment = getEnv("ENV", "development")

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		env, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return env
		}
	}
	return defaultValue
}

func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}
