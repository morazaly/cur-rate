package http_handler

import (
	"currency/config"
	"currency/service"
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	d *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{d: db}
}

func (h *Handler) Save(r *mux.Router, aconfig config.Config) error {

	r.HandleFunc("/currency/save/{date}", service.FetchFromApi(h.d)).Methods("GET")
	r.HandleFunc("/currency/{date}/{code}", service.GetFromApi(h.d)).Methods("GET")
	r.HandleFunc("/currency/{date}", service.GetFromApi(h.d)).Methods("GET")
	log.Println("Server started at ", aconfig.AppPort)

	log.Fatal(http.ListenAndServe(aconfig.AppPort, r))
	return nil
}
