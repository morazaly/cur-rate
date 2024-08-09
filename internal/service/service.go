package service

import (
	"context"
	"currency/internal/models"
	"currency/internal/repository"
	"encoding/json"
	"encoding/xml"
	"io"
	"log/slog"
	"net/http"
)

type Service struct {
	mySQLUserRepository *repository.MySQLUserRepository
	Log                 *slog.Logger
}

func New(
	mySQLUserRepository *repository.MySQLUserRepository,
	log *slog.Logger,
) *Service {
	return &Service{
		mySQLUserRepository: mySQLUserRepository,
		Log:                 log,
	}
}

func (s *Service) saveToDatabase(ctx context.Context, responseXml []byte, errCh chan error) {

	var rate models.Rates
	// Разбор XML
	if err := xml.Unmarshal(responseXml, &rate); err != nil {
		s.Log.Error("Failed to Unmarshal data: ", "err", err)
		errCh <- err
		return
	}

	for _, item := range rate.Items {
		item.Date = rate.Date

		count, err := s.mySQLUserRepository.Exists(ctx, &item, s.Log)
		if err != nil {
			s.Log.Error("Failed to query Exists: ", "err", err)
			errCh <- err
			return
		}
		if count == 0 {
			err := s.mySQLUserRepository.Insert(ctx, &item, s.Log)
			if err != nil {
				s.Log.Error("Failed to query Insert: ", "err", err)
				errCh <- err
				return
			}
		} else {
			err := s.mySQLUserRepository.Update(ctx, &item, s.Log)
			if err != nil {
				s.Log.Error("Failed to query Update: ", "err", err)
				errCh <- err
				return
			}
		}
	}
	errCh <- nil
}

func (s *Service) DownloadFromSource(ctx context.Context, url string, date string) models.Response {

	// URL of the public API
	apiURL := url + date

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
	defer close(errCh)
	go s.saveToDatabase(ctx, body, errCh)

	if err := <-errCh; err != nil {
		s.Log.Error("Failed to save data to database: ", "err", err)

		return models.Response{Success: false}
	}

	return models.Response{Success: true}

}

func (s *Service) GetSavedData(ctx context.Context, date string, code string) []byte {

	exists := code != ""

	var p []byte

	switch exists {

	case false:
		v, err := s.mySQLUserRepository.GetByDate(ctx, date, s.Log)
		if err != nil {

			s.Log.Error("Failed query: ", "err", err)

		}
		p, err = json.Marshal(v)
		if err != nil {

			s.Log.Error("Failed Marshal: ", "err", err)

		}

	case true:
		v, err := s.mySQLUserRepository.GetByDateCode(ctx, date, code, s.Log)
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
