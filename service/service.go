package service

import (
	"currency/models"
	"currency/repository"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func saveToDatabase(db *sql.DB, responseXml []byte, errCh chan error) {

	//errCh := make(chan error)

	//defer close(errCh)
	//defer db.Close()
	if db == nil {
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

		MySQLUserRepository := repository.MySQLUserRepository{Db: db}
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
	//	}()
	//return errCh
}

func FetchFromApi(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		date := vars["date"]

		// URL of the public API
		apiURL := "https://nationalbank.kz/rss/get_rates.cfm?fdate=" + date

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
		go saveToDatabase(db, body, errCh)

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

func GetFromApi(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		date := vars["date"]
		code, exists := vars["code"]
		if db == nil {
			return
		}
		MySQLUserRepository := repository.MySQLUserRepository{Db: db}
		switch exists {

		case false:

			v, err := MySQLUserRepository.GetByDate(date)

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
			v, err := MySQLUserRepository.GetByDateCode(date, code)
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
