package config

import (
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

type Config struct {
	HTTPServerAddress string
	DatabaseURL       string

	JWTSecret string
	JWTTTLMinutes int64
}

func New() *Config {
	secret := getEnv("JWT_SECRET", "change-me-in-prod")
	if secret == "change-me-in-prod" {
		log.Warn().Msg("JWT_SECRET is default; override it in production!")
	}

	ttl := int64(60)
	if v := os.Getenv("JWT_TTL_MINUTES"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil && n > 0 {
			ttl = n
		} else {
			log.Warn().Str("JWT_TTL_MINUTES", v).Msg("invalid JWT_TTL_MINUTES, fallback to 60")
		}
	}

	return &Config{
		HTTPServerAddress: getEnv("HTTP_SERVER_ADDRESS", ":8080"),
		DatabaseURL:       getEnv("DATABASE_URL", ""),
		JWTSecret:         secret,
		JWTTTLMinutes:     ttl,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}