package main

import (
	"currency/internal/app"
)

// @title Currency API
// @version 1.0
// @description This is a sample server for currency exchange rates.
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath

func main() {

	h := app.New()
	err := h.Run(app.Start)
	if err != nil {
		h.GetLogger().Error(err.Error())
	}
	h.GetLogger().Info("App is started")
}
