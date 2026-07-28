package config

import (
	"os"
	"strings"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	HTTPAddr    string
	CORSOrigins []string
}

func Load() Config {
	origins := os.Getenv("CORS_ORIGINS")
	if origins == "" {
		origins = "http://localhost:5173"
	}

	return Config{
		DatabaseURL: envOr("DATABASE_URL", "postgres://planeweave:planeweave@localhost:5432/planeweave?sslmode=disable"),
		JWTSecret:   envOr("JWT_SECRET", "dev-secret-change-me"),
		HTTPAddr:    envOr("HTTP_ADDR", ":8080"),
		CORSOrigins: strings.Split(origins, ","),
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
