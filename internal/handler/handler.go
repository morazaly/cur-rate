package handler

import (
	"context"
	"currency/internal/config"
	"currency/internal/service"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	aconfig config.Config
	service *service.Service
}

func NewHandler(aconfig config.Config, service *service.Service) *Handler {
	return &Handler{aconfig: aconfig,
		service: service}
}

func (h *Handler) StartHandler(ch chan error) {

	r := mux.NewRouter()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r.HandleFunc("/currency/save/{date}", h.downloadFromSource(ctx)).Methods("GET")
	r.HandleFunc("/currency/{date}/{code}", h.getSavedData(ctx)).Methods("GET")
	r.HandleFunc("/currency/{date}", h.getSavedData(ctx)).Methods("GET")
	ch <- http.ListenAndServe(h.aconfig.AppPort, r)

}

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
