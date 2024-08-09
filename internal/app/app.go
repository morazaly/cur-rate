package app

import (
	"currency/internal/config"
	"currency/internal/db"
	"currency/internal/handler"
	"currency/internal/logger"
	"currency/internal/repository"
	"currency/internal/service"
	"database/sql"
	"log/slog"

	"os"
	"os/signal"
	"syscall"
)

type app struct {
	log *slog.Logger
}

type App interface {
	GetLogger() *slog.Logger
	Start() error
}

func New() App {

	logger := logger.InitLogger()
	return &app{log: logger}
}

func (h *app) Start() error {

	aconfig := *config.NewConfig()

	errChan := make(chan error)
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	// Инициализация БД
	db := initDb(aconfig)
	repository := repository.New(db)
	s := service.New(repository, h.GetLogger())
	handler := handler.NewHandler(aconfig, s)
	h.GetLogger().Info("Server started at ", "Apport ", aconfig.AppPort)

	go handler.StartHandler(errChan)

	select {
	case err := <-errChan:
		h.GetLogger().Error("gracefully shutdown error", "error", err.Error())
	case stop := <-stopChan:
		h.GetLogger().Error("app is finished", "signal", stop.String())
	}

	return nil
}

func initDb(aconfig config.Config) *sql.DB {

	// Подключение к базе данных MySQL
	adb := db.NewDb(&aconfig)
	return adb

}

func (a app) GetLogger() *slog.Logger {
	return a.log
}
