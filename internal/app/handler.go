package app

import (
	"context"
	"currency/internal/models"
	"currency/internal/service"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	aconfig models.Config
	service *service.Service
}

func NewHandler(aconfig models.Config, service *service.Service) *Handler {
	return &Handler{aconfig: aconfig,
		service: service}
}

func (h *Handler) StartHandler() error {
	r := mux.NewRouter()
	r.HandleFunc("/currency/save/{date}", fetchApi(h.aconfig, h.service, "DownloadFromSource")).Methods("GET")
	r.HandleFunc("/currency/{date}/{code}", fetchApi(h.aconfig, h.service, "GetSavedData")).Methods("GET")
	r.HandleFunc("/currency/{date}", fetchApi(h.aconfig, h.service, "GetSavedData")).Methods("GET")
	return http.ListenAndServe(h.aconfig.AppPort, r)
}

func fetchApi(aconfig models.Config, s *service.Service, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		date := vars["date"]
		code := vars["code"]
		var bytes []byte
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		w.Header().Set("Content-Type", "application/json")
		switch method {
		case "GetSavedData":
			bytes = s.GetSavedData(ctx, date, code)
			w.Write(bytes)
		case "DownloadFromSource":
			response := s.DownloadFromSource(ctx, aconfig, date)
			if err := json.NewEncoder(w).Encode(response); err != nil {
				s.Log.Error("Failed Marshal: ", "err", err)
			}

		}
	}
}
