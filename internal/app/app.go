package app

import (
	"currency/internal/db"
	"currency/internal/models"
	"currency/internal/repository"
	"currency/internal/service"
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Start(r *mux.Router) error {
	// Загрузка конфигурации из файла config.json
	aconfig := *models.NewConfig()
	// Инициализация БД
	db := initDb(aconfig)
	r.HandleFunc("/currency/save/{date}", service.New(repository.MySQLUserRepository{Db: db}).FetchFromApi(aconfig)).Methods("GET")
	r.HandleFunc("/currency/{date}/{code}", service.New(repository.MySQLUserRepository{Db: db}).GetFromApi()).Methods("GET")
	r.HandleFunc("/currency/{date}", service.New(repository.MySQLUserRepository{Db: db}).GetFromApi()).Methods("GET")
	log.Println("Server started at ", aconfig.AppPort)

	log.Fatal(http.ListenAndServe(aconfig.AppPort, r))
	return nil
}

func initDb(aconfig models.Config) *sql.DB {

	// Подключение к базе данных MySQL
	adb := db.NewDb(&aconfig)
	return adb

}
