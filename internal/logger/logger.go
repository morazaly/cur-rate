package logger

import (
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func InitLogger() *slog.Logger {
	if err := godotenv.Load("..\\.env"); err != nil {
		panic("Error loading .env file")
	}
	logLevel := os.Getenv("LOG_LEVEL")
	logFormat := os.Getenv("LOG_FORMAT")
	if logLevel == "" {
		logLevel = "info"
	}
	if logFormat == "" {
		logFormat = "json"
	}

	var level slog.Level
	switch strings.ToLower(logLevel) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	handlerOptions := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if strings.ToLower(logFormat) == "json" {
		handler = slog.NewJSONHandler(os.Stdout, handlerOptions)
	} else {
		handler = slog.NewTextHandler(os.Stdout, handlerOptions)
	}

	return slog.New(handler)
}
