package app

import (
	"currency/internal/db"
	"currency/internal/models"
	"currency/internal/repository"
	"currency/internal/service"
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type app struct {
	log *slog.Logger
}

type App interface {
	GetLogger() *slog.Logger
	Start(r *mux.Router) error
}

func New() App {

	handler := slog.NewJSONHandler(os.Stdout, nil)

	logger := slog.New(handler)

	return &app{log: logger}
}

func (h *app) Start(r *mux.Router) error {
	// Загрузка конфигурации из файла config.json
	aconfig := *models.NewConfig()
	// Инициализация БД
	db := initDb(aconfig)
	r.HandleFunc("/currency/save/{date}", service.New(repository.MySQLUserRepository{Db: db}, h.GetLogger()).FetchFromApi(aconfig)).Methods("GET")
	r.HandleFunc("/currency/{date}/{code}", service.New(repository.MySQLUserRepository{Db: db}, h.GetLogger()).GetFromApi()).Methods("GET")
	r.HandleFunc("/currency/{date}", service.New(repository.MySQLUserRepository{Db: db}, h.GetLogger()).GetFromApi()).Methods("GET")
	h.GetLogger().Info("Server started at ", "Apport ", aconfig.AppPort)
	//log.Fatal(http.ListenAndServe(aconfig.AppPort, r))
	h.GetLogger().Error(http.ListenAndServe(aconfig.AppPort, r).Error())
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
