package service

import (
	"currency/internal/models"
	"currency/internal/repository"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Service struct {
	MySQLUserRepository repository.MySQLUserRepository
}

func New(
	MySQLUserRepository repository.MySQLUserRepository,
) *Service {
	return &Service{
		MySQLUserRepository: MySQLUserRepository,
	}
}

func saveToDatabase(MySQLUserRepository repository.MySQLUserRepository, responseXml []byte, errCh chan error) {

	if MySQLUserRepository.Db == nil {
		errCh <- fmt.Errorf("database connection is nil")
		return
	}
	var rate models.Rates
	// Разбор XML
	if err := xml.Unmarshal(responseXml, &rate); err != nil {
		errCh <- err
		return
	}

	for _, item := range rate.Items {
		item.Date = rate.Date
		count, err := MySQLUserRepository.Exists(&item)
		if err != nil {
			errCh <- err
			return
		}
		if count == 0 {
			err := MySQLUserRepository.Insert(&item)
			if err != nil {
				errCh <- err
				return
			}
		} else {
			err := MySQLUserRepository.Update(&item)
			if err != nil {
				errCh <- err
				return
			}
		}
	}
	errCh <- nil
}

func (s *Service) FetchFromApi(aconfig models.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		date := vars["date"]

		// URL of the public API
		apiURL := aconfig.ApiURL + date

		// Make the HTTP GET request
		resp, err := http.Get(apiURL)
		if err != nil {
			log.Printf("Failed to fetch data: %v", err)
			http.Error(w, "Failed to fetch data"+err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response: %v", err)
			http.Error(w, "Failed to read response", http.StatusInternalServerError)
			return
		}
		// Сохранение в базу данных
		errCh := make(chan error)
		//defer close(errCh)
		go saveToDatabase(s.MySQLUserRepository, body, errCh)

		if err := <-errCh; err != nil {
			log.Printf("Failed to save data to database: %v", err)
			http.Error(w, "Failed to save data to database"+err.Error(), http.StatusInternalServerError)
			return
		}

		// Write the JSON response
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(models.Response{Success: true}); err != nil {
			log.Printf("Failed Marshal: %v", err)
			http.Error(w, "Failed Marshal", http.StatusInternalServerError)
			return
		}
	}
}

func (s *Service) GetFromApi() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		date := vars["date"]
		code, exists := vars["code"]
		if s.MySQLUserRepository.Db == nil {
			return
		}
		switch exists {

		case false:
			v, err := s.MySQLUserRepository.GetByDate(date)
			if err != nil {
				return
			}
			p, err := json.Marshal(v)
			if err != nil {
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(p)

		case true:
			v, err := s.MySQLUserRepository.GetByDateCode(date, code)
			if err != nil {
				return
			}
			p, err := json.Marshal(v)
			if err != nil {
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(p)
		}

	}
}
