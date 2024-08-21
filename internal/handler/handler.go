package handler

import (
	"context"
	"currency/internal/config"
	"currency/internal/service"
	"encoding/json"
	"net/http"

	_ "currency/docs"

	"github.com/gorilla/mux"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Handler struct {
	aconfig config.Config
	service *service.Service
}

func NewHandler(aconfig config.Config, service *service.Service) *Handler {
	return &Handler{aconfig: aconfig,
		service: service}
}

func (h *Handler) StartHandler(ctx context.Context, ch chan error) {

	r := mux.NewRouter()

	go func() {
		ch <- h.service.Metrics.Start(ctx)
	}()
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	r.HandleFunc("/currency/save/{date}", h.downloadFromSource(ctx)).Methods("GET")
	r.HandleFunc("/currency/{date}/{code}", h.getSavedData(ctx)).Methods("GET")
	r.HandleFunc("/currency/{date}", h.getSavedData(ctx)).Methods("GET")
	ch <- http.ListenAndServe(h.aconfig.AppPort, r)

}

// downloadFromSource example
// @Summary Save currency by date
// @Description Save currency by date
// @Tags downloadFromSource
// @Accept  json
// @Produce  json
// @Param date path string true "Date"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /currency/save/{date} [get]
func (h *Handler) downloadFromSource(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		date := vars["date"]
		url := h.aconfig.ApiURL
		w.Header().Set("Content-Type", "application/json")

		response := h.service.DownloadFromSource(ctx, url, date)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			h.service.Log.Error("Failed Marshal: ", "err", err)
		}

	}
}

// getSavedData example
// @Summary Get currency by date and code
// @Description  Get currency exchange rate by date and optionally by code
// @Tags getSavedData
// @Accept  json
// @Produce  json
// @Param date path string true "Date"
// @Param code path string false "Code"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /currency/{date}/{code} [get]
// @Router /currency/{date} [get]
func (h *Handler) getSavedData(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		date := vars["date"]
		code := vars["code"]
		var bytes []byte
		w.Header().Set("Content-Type", "application/json")

		bytes = h.service.GetSavedData(ctx, date, code)
		w.Write(bytes)

	}
}
