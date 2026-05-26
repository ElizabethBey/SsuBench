package config

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	JWTSecret       string
	HTTPAddr        string
	LogLevel        slog.Level
	ShutdownTimeout time.Duration
}

func FromEnv() Config {
	return Config{
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPassword:      getEnv("DB_PASSWORD", "postgres"),
		DBName:          getEnv("DB_NAME", "ssubench"),
		JWTSecret:       getEnv("JWT_SECRET", "default_secret_key"),
		HTTPAddr:        getEnv("HTTP_ADDR", ":8080"),
		LogLevel:        parseLogLevel(getEnv("LOG_LEVEL", "INFO")),
		ShutdownTimeout: parseDuration(getEnv("SHUTDOWN_TIMEOUT", "5s"), 5*time.Second),
	}
}

func getEnv(key string, def string) string {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return def
	}
	return val
}

func parseLogLevel(s string) slog.Level {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN", "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func parseDuration(s string, def time.Duration) time.Duration {
	d, err := time.ParseDuration(strings.TrimSpace(s))
	if err != nil {
		return def
	}
	return d
}

// Иногда удобно иметь int env-парсер под будущие настройки.
func parseInt(s string, def int) int {
	v, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return def
	}
	return v
}
