package main

import (
	"currency/internal/app"

	"github.com/gorilla/mux"
)

func main() {

	// Создание маршрута и запуск веб-сервера
	r := mux.NewRouter()
	h := app.New()
	err := h.Start(r)
	if err != nil {
		h.GetLogger().Error(err.Error())
	}
	h.GetLogger().Info("App is started")
}
