package service

import (
	"currency/models"
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

	// Пример запроса на вставку данных в таблицу

	// если есть запись COD - A_DATE, сделать апдейт или пропускать

	for _, item := range rate.Items {
		//	for  int i:= 0; i++; i< len(Items) - 1 {
		var count int

		err := db.QueryRow("SELECT COUNT(*) FROM  r_currency  where code = ? AND A_DATE = str_to_date(?,'%d.%m.%Y')", item.Title, rate.Date).Scan(&count)

		if err != nil {
			errCh <- err
			return
		}

		if count == 0 {
			_, err := db.Exec("INSERT INTO r_currency (TITLE, CODE, VALUE, A_DATE) VALUES (?, ?, ?, str_to_date(?,'%d.%m.%Y'))", item.Fullname, item.Title, item.Description, rate.Date)

			if err != nil {
				errCh <- err
				return
			}

		} else {
			_, err := db.Exec("UPDATE r_currency SET VALUE = ? WHERE  CODE = ? AND  A_DATE = str_to_date(?,'%d.%m.%Y')", item.Description, item.Title, rate.Date)

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
		go saveToDatabase(db, body, errCh) // 10

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
		// w.Write(data)
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

		switch exists {

		case false:

			sql, err := db.Query("SELECT * FROM  r_currency  where  A_DATE = str_to_date(?,'%d.%m.%Y')", date)

			if err != nil {

				return
			}

			defer sql.Close()

			var v struct {
				Data []models.ResponseItem `json:"data"`
			}

			for sql.Next() {
				var responseItem models.ResponseItem

				if err := sql.Scan(&responseItem.Id, &responseItem.Title, &responseItem.Code, &responseItem.Value, &responseItem.Adate); err != nil {
					// handle error
					return
				}
				v.Data = append(v.Data, responseItem)
			}
			p, err := json.Marshal(v)

			if err != nil {
				// handle error
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(p)

		case true:

			sql, err := db.Query("SELECT * FROM  r_currency  where code = ? AND A_DATE = str_to_date(?,'%d.%m.%Y')", code, date)

			if err != nil {

				return
			}
			defer sql.Close()

			var v struct {
				Data []models.ResponseItem `json:"data"`
			}

			for sql.Next() {
				var responseItem models.ResponseItem

				if err := sql.Scan(&responseItem.Id, &responseItem.Title, &responseItem.Code, &responseItem.Value, &responseItem.Adate); err != nil {
					// handle error
					return
				}
				v.Data = append(v.Data, responseItem)
			}

			p, err := json.Marshal(v)

			if err != nil {
				// handle error
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(p)
		}

	}
}
