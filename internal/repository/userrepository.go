package repository

import (
	"context"
	"currency/internal/metrics"
	"currency/internal/models"
	"database/sql"
	"log/slog"
	"time"
)

type MySQLUserRepository struct {
	// DB connection and other fields
	db *sql.DB
}

func New(
	db *sql.DB,
) *MySQLUserRepository {
	return &MySQLUserRepository{
		db: db,
	}
}

type UserRepository interface {
	GetByDate(Date string) ([]models.ResponseItem, error)
	GetByDateCode(Date string, Code string) ([]models.ResponseItem, error)
	Exists(user *models.Item) (int, error)
	Update(user *models.Item) error
	Insert(user *models.Item) error
}

const getByDateQuery = "SELECT * FROM  r_currency  where  A_DATE = str_to_date(?,'%d.%m.%Y')"

func (repo *MySQLUserRepository) GetByDate(ctx context.Context, date string, logger *slog.Logger) ([]models.ResponseItem, error) {
	// Implementation
	const op = "GetByDate"
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	now := time.Now()
	defer func() {
		go metrics.DurPGProcessed.WithLabelValues(op).Observe(time.Since(now).Seconds())
	}()

	sql, err := repo.db.QueryContext(ctxT, getByDateQuery, date)

	if err != nil {
		logger.Error("Failed to query Exists: ", "err", err)
		return nil, err
	}

	defer sql.Close()

	var Data []models.ResponseItem

	for sql.Next() {
		var responseItem models.ResponseItem

		if err := sql.Scan(&responseItem.Id, &responseItem.Title, &responseItem.Code, &responseItem.Value, &responseItem.Adate); err != nil {
			logger.Error("Failed to query Exists: ", "err", err)
			return nil, err
		}
		Data = append(Data, responseItem)
	}

	return Data, nil
}

const getByDateCodeQuery = "SELECT * FROM  r_currency  where code = ? AND A_DATE = str_to_date(?,'%d.%m.%Y')"

func (repo *MySQLUserRepository) GetByDateCode(ctx context.Context, date string, code string, logger *slog.Logger) ([]models.ResponseItem, error) {
	// Implementation
	const op = "GetByDateCode"
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	now := time.Now()
	defer func() {
		go metrics.DurPGProcessed.WithLabelValues(op).Observe(time.Since(now).Seconds())
	}()
	sql, err := repo.db.QueryContext(ctxT, getByDateCodeQuery, code, date)

	if err != nil {
		logger.Error("Failed to query Exists: ", "err", err)
		return nil, err
	}
	defer sql.Close()

	var v struct {
		Data []models.ResponseItem `json:"data"`
	}

	for sql.Next() {
		var responseItem models.ResponseItem

		if err := sql.Scan(&responseItem.Id, &responseItem.Title, &responseItem.Code, &responseItem.Value, &responseItem.Adate); err != nil {
			logger.Error("Failed to query Exists: ", "err", err)
			return nil, err
		}
		v.Data = append(v.Data, responseItem)
	}

	return v.Data, nil
}

const existsQuery = "SELECT COUNT(*) FROM  r_currency  where code = ? AND A_DATE = str_to_date(?,'%d.%m.%Y')"

const updateQuery = "UPDATE r_currency SET VALUE = ? WHERE  CODE = ? AND  A_DATE = str_to_date(?,'%d.%m.%Y')"

const insertQuery = "INSERT INTO r_currency (TITLE, CODE, VALUE, A_DATE) VALUES (?, ?, ?, str_to_date(?,'%d.%m.%Y'))"

func (repo *MySQLUserRepository) Append(ctx context.Context, user *models.Item, logger *slog.Logger) error {
	// Implementation
	const op = "AppendQuery"
	var count int
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	now := time.Now()
	defer func() {
		go metrics.DurPGProcessed.WithLabelValues(op).Observe(time.Since(now).Seconds())
	}()

	err := repo.db.QueryRowContext(ctxT, existsQuery, user.Title, user.Date).Scan(&count)

	if err != nil {
		logger.Error("Failed to query Exists: ", "err", err)
		return err

	}

	if count == 0 {
		_, err := repo.db.ExecContext(ctx, insertQuery, user.Fullname, user.Title, user.Description, user.Date)

		if err != nil {
			logger.Error("Failed to query Exists: ", "err", err)
			return err
		}

	} else {

		_, err := repo.db.ExecContext(ctx, updateQuery, user.Description, user.Title, user.Date)
		if err != nil {
			logger.Error("Failed to query Exists: ", "err", err)
			return err
		}

	}
	return nil
}
