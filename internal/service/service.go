package service

import (
	"context"
	"currency/internal/models"
	"currency/internal/repository"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type Service struct {
	MySQLUserRepository repository.MySQLUserRepository
	Log                 *slog.Logger
}

func New(
	MySQLUserRepository repository.MySQLUserRepository,
	Log *slog.Logger,
) *Service {
	return &Service{
		MySQLUserRepository: MySQLUserRepository,
		Log:                 Log,
	}
}

func (s *Service) saveToDatabase(ctx context.Context, responseXml []byte, errCh chan error) {

	MySQLUserRepository := s.MySQLUserRepository

	if MySQLUserRepository.Db == nil {
		errCh <- fmt.Errorf("database connection is nil")
		return
	}
	var rate models.Rates
	// Разбор XML
	if err := xml.Unmarshal(responseXml, &rate); err != nil {
		s.Log.Error("Failed to Unmarshal data: ", "err", err)
		errCh <- err
		return
	}

	for _, item := range rate.Items {
		item.Date = rate.Date
		count, err := MySQLUserRepository.Exists(ctx, &item)
		if err != nil {
			s.Log.Error("Failed to query Exists: ", "err", err)
			errCh <- err
			return
		}
		if count == 0 {
			err := MySQLUserRepository.Insert(ctx, &item)
			if err != nil {
				s.Log.Error("Failed to query Insert: ", "err", err)
				errCh <- err
				return
			}
		} else {
			err := MySQLUserRepository.Update(ctx, &item)
			if err != nil {
				s.Log.Error("Failed to query Update: ", "err", err)
				errCh <- err
				return
			}
		}
	}
	errCh <- nil
}

func (s *Service) DownloadFromSource(ctx context.Context, aconfig models.Config, date string) models.Response {

	// URL of the public API
	apiURL := aconfig.ApiURL + date

	// Make the HTTP GET request
	resp, err := http.Get(apiURL)
	if err != nil {
		s.Log.Error("Failed to fetch data: ", "err", err)
		return models.Response{Success: false}
	}
	defer resp.Body.Close()
	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.Log.Error("Failed to read response: ", "err", err)
		return models.Response{Success: false}
	}
	// Сохранение в базу данных
	errCh := make(chan error)
	//defer close(errCh)
	go s.saveToDatabase(ctx, body, errCh)

	if err := <-errCh; err != nil {
		s.Log.Error("Failed to save data to database: ", "err", err)
		return models.Response{Success: false}
	}

	return models.Response{Success: true}

}

func (s *Service) GetSavedData(ctx context.Context, date string, code string) []byte {

	exists := code != ""
	if s.MySQLUserRepository.Db == nil {
		s.Log.Error("db connection is null")
	}
	var p []byte

	switch exists {

	case false:
		v, err := s.MySQLUserRepository.GetByDate(ctx, date)
		if err != nil {
			s.Log.Error("Failed query: ", "err", err)

		}
		p, err = json.Marshal(v)
		if err != nil {
			s.Log.Error("Failed Marshal: ", "err", err)

		}

	case true:
		v, err := s.MySQLUserRepository.GetByDateCode(ctx, date, code)
		if err != nil {
			s.Log.Error("Failed query: ", "err", err)

		}
		p, err = json.Marshal(v)
		if err != nil {
			s.Log.Error("Failed Marshal: ", "err", err)

		}

	}
	return p
}
