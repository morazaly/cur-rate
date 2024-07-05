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
	// Загрузка конфигурации из файла config.json
	aconfig = *config.NewConfig()

	// Подключение к базе данных MySQL
	adb = db.NewDb(&aconfig)

}

func main() {

	// Создание маршрута и запуск веб-сервера
	r := mux.NewRouter()
	h := http_handler.NewHandler(adb)
	h.Save(r, aconfig)

}
