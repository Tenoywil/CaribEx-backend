package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	// Environment
	AppEnv string `mapstructure:"ENV"`

	// Server Configuration
	ServerPort            string `mapstructure:"PORT"`
	ServerHost            string `mapstructure:"HOST"`
	ServerReadTimeout     string `mapstructure:"SERVER_READ_TIMEOUT"`
	ServerWriteTimeout    string `mapstructure:"SERVER_WRITE_TIMEOUT"`
	ServerShutdownTimeout string `mapstructure:"SERVER_SHUTDOWN_TIMEOUT"`
	AllowedOrigins        string `mapstructure:"ALLOWED_ORIGINS"`

	// Database Configuration
	DBConnectionString string `mapstructure:"DB_CONNECTION_STRING"`
	DBMaxConnections   int    `mapstructure:"DB_MAX_CONNECTIONS"`
	DBMaxIdleTime      string `mapstructure:"DB_MAX_IDLE_TIME"`
	DBMaxConnLifetime  string `mapstructure:"DB_MAX_CONN_LIFETIME"`

	// Redis Configuration
	RedisConnectionString string `mapstructure:"REDIS_CONNECTION_STRING"`
	RedisHost             string `mapstructure:"REDIS_HOST"`
	RedisPort             int    `mapstructure:"REDIS_PORT"`
	RedisPassword         string `mapstructure:"REDIS_PASSWORD"`
	RedisDB               int    `mapstructure:"REDIS_DB"`

	// Authentication Configuration
	SessionSecret   string `mapstructure:"SESSION_SECRET"`
	SessionDuration string `mapstructure:"SESSION_DURATION"`
	JWTSecret       string `mapstructure:"JWT_SECRET"`
	JWTExpiration   string `mapstructure:"JWT_EXPIRATION"`
	SIWEDomain      string `mapstructure:"SIWE_DOMAIN"`

	// Cache Configuration
	CacheEnableL1  bool   `mapstructure:"CACHE_ENABLE_L1"`
	CacheEnableL2  bool   `mapstructure:"CACHE_ENABLE_L2"`
	CacheL1MaxSize int64  `mapstructure:"CACHE_L1_MAX_SIZE"`
	CacheL1TTL     string `mapstructure:"CACHE_L1_TTL"`
	CacheL2TTL     string `mapstructure:"CACHE_L2_TTL"`

	// Blockchain Configuration
	RPCURL string `mapstructure:"RPC_URL"`

	// Parsed values
	AllowedOriginsSlice []string
}

// Load loads configuration from environment variables
func Load() *Config {
	cfg := &Config{}

	// Try to read from .env first (dev/local). If not found and running in production, fall back to OS env.
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		if os.Getenv("ENV") == "production" {
			loadEnvFromOS(cfg)
			return cfg
		}
		log.Fatal("Can't find the file .env : ", err)
	}

	if err := viper.Unmarshal(cfg); err != nil {
		if os.Getenv("ENV") == "production" {
			loadEnvFromOS(cfg)
			return cfg
		}
		log.Fatal("Environment can't be loaded: ", err)
	}

	// If the loaded configuration indicates production, prefer OS env (e.g., when .env exists but prod uses env vars)
	switch cfg.AppEnv {
	case "production":
		loadEnvFromOS(cfg)
	case "development":
		log.Println("The App is running in development env")
	}

	// Parse allowed origins into slice
	cfg.AllowedOriginsSlice = allowedOriginSlice(cfg.AllowedOrigins)
	log.Printf("[CONFIG] Loaded ALLOWED_ORIGINS: %s", cfg.AllowedOrigins)
	log.Printf("[CONFIG] Parsed AllowedOriginsSlice: %v", cfg.AllowedOriginsSlice)

	return cfg
}

func loadEnvFromOS(cfg *Config) {
	cfg.AppEnv = os.Getenv("ENV")

	// Server Configuration
	cfg.ServerPort = os.Getenv("PORT")
	cfg.ServerHost = os.Getenv("HOST")
	cfg.ServerReadTimeout = os.Getenv("SERVER_READ_TIMEOUT")
	cfg.ServerWriteTimeout = os.Getenv("SERVER_WRITE_TIMEOUT")
	cfg.ServerShutdownTimeout = os.Getenv("SERVER_SHUTDOWN_TIMEOUT")
	cfg.AllowedOrigins = os.Getenv("ALLOWED_ORIGINS")

	// Database Configuration
	cfg.DBConnectionString = os.Getenv("DB_CONNECTION_STRING")
	cfg.DBMaxConnections = getenvInt("DB_MAX_CONNECTIONS")
	cfg.DBMaxIdleTime = os.Getenv("DB_MAX_IDLE_TIME")
	cfg.DBMaxConnLifetime = os.Getenv("DB_MAX_CONN_LIFETIME")

	// Redis Configuration
	cfg.RedisConnectionString = os.Getenv("REDIS_CONNECTION_STRING")
	cfg.RedisHost = os.Getenv("REDIS_HOST")
	cfg.RedisPort = getenvInt("REDIS_PORT")
	cfg.RedisPassword = os.Getenv("REDIS_PASSWORD")
	cfg.RedisDB = getenvInt("REDIS_DB")

	// Authentication Configuration
	cfg.SessionSecret = os.Getenv("SESSION_SECRET")
	cfg.SessionDuration = os.Getenv("SESSION_DURATION")
	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	cfg.JWTExpiration = os.Getenv("JWT_EXPIRATION")
	cfg.SIWEDomain = os.Getenv("SIWE_DOMAIN")

	// Cache Configuration
	cfg.CacheEnableL1 = getenvBool("CACHE_ENABLE_L1")
	cfg.CacheEnableL2 = getenvBool("CACHE_ENABLE_L2")
	cfg.CacheL1MaxSize = getenvInt64("CACHE_L1_MAX_SIZE")
	cfg.CacheL1TTL = os.Getenv("CACHE_L1_TTL")
	cfg.CacheL2TTL = os.Getenv("CACHE_L2_TTL")

	// Blockchain Configuration
	cfg.RPCURL = os.Getenv("RPC_URL")

	// Parse allowed origins into slice
	cfg.AllowedOriginsSlice = allowedOriginSlice(cfg.AllowedOrigins)
}

func getenvInt(key string) int {
	v := os.Getenv(key)
	if v == "" {
		return 0
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return i
}

func getenvInt64(key string) int64 {
	v := os.Getenv(key)
	if v == "" {
		return 0
	}
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

func getenvBool(key string) bool {
	v := os.Getenv(key)
	if v == "" {
		return false
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false
	}
	return b
}

func allowedOriginSlice(origins string) []string {
	if origins == "" {
		log.Println("[CONFIG] WARNING: ALLOWED_ORIGINS is empty! CORS will not work properly.")
		return []string{}
	}
	var result []string
	for _, origin := range strings.Split(origins, ",") {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
