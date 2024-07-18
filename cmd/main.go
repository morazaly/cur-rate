package main

import (
	"currency/internal/app"

	"github.com/gorilla/mux"
)

func main() {

	// Создание маршрута и запуск веб-сервера
	r := mux.NewRouter()
	h := app.New()
	h.Start(r)

}
