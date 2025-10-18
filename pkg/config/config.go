package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Auth     AuthConfig
	Cache    CacheConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port            string
	Host            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	AllowedOrigins  []string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	MaxConnections  int
	MaxIdleTime     time.Duration
	MaxConnLifetime time.Duration
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	SessionSecret   string
	SessionDuration time.Duration
	JWTSecret       string
	JWTExpiration   time.Duration
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	L1MaxSize     int64
	L1TTL         time.Duration
	L2TTL         time.Duration
	EnableL1      bool
	EnableL2      bool
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:            getEnv("PORT", "8080"),
			Host:            getEnv("HOST", "0.0.0.0"),
			ReadTimeout:     getDurationEnv("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout:    getDurationEnv("SERVER_WRITE_TIMEOUT", 10*time.Second),
			ShutdownTimeout: getDurationEnv("SERVER_SHUTDOWN_TIMEOUT", 30*time.Second),
			AllowedOrigins:  []string{getEnv("ALLOWED_ORIGIN", "http://localhost:3000")},
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getIntEnv("DB_PORT", 5432),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "postgres"),
			Database:        getEnv("DB_NAME", "caribx"),
			MaxConnections:  getIntEnv("DB_MAX_CONNECTIONS", 25),
			MaxIdleTime:     getDurationEnv("DB_MAX_IDLE_TIME", 15*time.Minute),
			MaxConnLifetime: getDurationEnv("DB_MAX_CONN_LIFETIME", 1*time.Hour),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getIntEnv("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getIntEnv("REDIS_DB", 0),
		},
		Auth: AuthConfig{
			SessionSecret:   getEnv("SESSION_SECRET", "change-me-in-production"),
			SessionDuration: getDurationEnv("SESSION_DURATION", 24*time.Hour),
			JWTSecret:       getEnv("JWT_SECRET", "change-me-in-production"),
			JWTExpiration:   getDurationEnv("JWT_EXPIRATION", 1*time.Hour),
		},
		Cache: CacheConfig{
			L1MaxSize:     getInt64Env("CACHE_L1_MAX_SIZE", 100*1024*1024), // 100MB
			L1TTL:         getDurationEnv("CACHE_L1_TTL", 5*time.Minute),
			L2TTL:         getDurationEnv("CACHE_L2_TTL", 15*time.Minute),
			EnableL1:      getBoolEnv("CACHE_ENABLE_L1", true),
			EnableL2:      getBoolEnv("CACHE_ENABLE_L2", true),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}
	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Auth.SessionSecret == "change-me-in-production" && os.Getenv("ENV") == "production" {
		return fmt.Errorf("session secret must be set in production")
	}
	return nil
}

// Helper functions to get environment variables with defaults
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getInt64Env(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
