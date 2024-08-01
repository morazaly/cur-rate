package main

import (
	"currency/internal/app"
)

func main() {

	h := app.New()
	err := h.Start()
	if err != nil {
		h.GetLogger().Error(err.Error())
	}
	h.GetLogger().Info("App is started")
}
