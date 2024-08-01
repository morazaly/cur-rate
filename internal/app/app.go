package app

import (
	"currency/internal/db"
	"currency/internal/models"
	"currency/internal/repository"
	"currency/internal/service"
	"database/sql"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type app struct {
	log *slog.Logger
}

type App interface {
	GetLogger() *slog.Logger
	Start() error
}

func New() App {

	logger := initLogger()
	return &app{log: logger}
}

func (h *app) Start() error {

	aconfig := *models.NewConfig()

	// Инициализация БД
	db := initDb(aconfig)
	s := service.New(repository.MySQLUserRepository{Db: db}, h.GetLogger())
	handler := NewHandler(aconfig, s)
	h.GetLogger().Info("Server started at ", "Apport ", aconfig.AppPort)
	h.GetLogger().Error(handler.StartHandler().Error())
	return nil
}

func initDb(aconfig models.Config) *sql.DB {

	// Подключение к базе данных MySQL
	adb := db.NewDb(&aconfig)
	return adb

}

func (a app) GetLogger() *slog.Logger {
	return a.log
}

func initLogger() *slog.Logger {
	if err := godotenv.Load(); err != nil {
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
