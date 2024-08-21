package app

import (
	"context"
	"currency/internal/config"
	"currency/internal/db"
	"currency/internal/handler"
	"currency/internal/logger"
	"currency/internal/metrics"
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
	//Start(errChan chan error) error
	//Run() error
	Run(Start func(h *app, handler *handler.Handler) error) error
}

func (h *app) Run(Start func(h *app, handler *handler.Handler) error) error {

	aconfig := *config.NewConfig()

	// Инициализация БД
	db := initDb(aconfig)
	repository := repository.New(db)
	// Инициализация Метрик
	metrics := initMetrics(h.GetLogger(), aconfig)
	// Инициализация Service
	s := service.New(repository, h.GetLogger(), metrics)
	// Инициализация Handler
	handler := handler.NewHandler(aconfig, s)

	h.GetLogger().Info("Server started at ", "Apport ", aconfig.AppPort)

	return Start(h, handler)
}

func New() App {

	logger := logger.InitLogger()
	return &app{log: logger}
}

func Start(h *app, handler *handler.Handler) error {

	errChan := make(chan error)
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go handler.StartHandler(ctx, errChan)
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

func initMetrics(log *slog.Logger, cfg config.Config) *metrics.Metrics {
	var (
		metricsPort string
	)

	metricsPort = cfg.MetricPort
	return metrics.NewMetrics(log, metricsPort)
}
