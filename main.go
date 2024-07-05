package main

import (
	"currency/config"
	"currency/db"
	http_handler "currency/handler"
	"database/sql"

	"github.com/gorilla/mux"
)

var adb *sql.DB
var aconfig config.Config

func init() {
	//var err error

	// Загрузка конфигурации из файла config.json
	aconfig = *config.NewConfig()

	// Подключение к базе данных MySQL
	adb = db.NewDb(&aconfig)

}

func main() {

	// Создание маршрута и запуск веб-сервера

	/*r := mux.NewRouter()
	r.HandleFunc("/currency/save/{date}", service.FetchFromApi(db)).Methods("GET")
	r.HandleFunc("/currency/{date}/{code}", service.GetFromApi(db)).Methods("GET")
	r.HandleFunc("/currency/{date}", service.GetFromApi(db)).Methods("GET")

	log.Println("Server started at ", aconfig.AppPort)
	log.Fatal(http.ListenAndServe(aconfig.AppPort, r))*/
	r := mux.NewRouter()
	h := http_handler.NewHandler(adb)
	h.Save(r, aconfig)

}
