package repository

import (
	"context"
	"currency/internal/models"
	"database/sql"
)

type MySQLUserRepository struct {
	// DB connection and other fields
	Db *sql.DB
}

type UserRepository interface {
	GetByDate(Date string) ([]models.ResponseItem, error)
	GetByDateCode(Date string, Code string) ([]models.ResponseItem, error)
	Exists(user *models.Item) (int, error)
	Update(user *models.Item) error
	Insert(user *models.Item) error
}

const getByDateQuery = "SELECT * FROM  r_currency  where  A_DATE = str_to_date(?,'%d.%m.%Y')"

func (repo *MySQLUserRepository) GetByDate(ctx context.Context, date string) ([]models.ResponseItem, error) {
	// Implementation

	sql, err := repo.Db.QueryContext(ctx, getByDateQuery, date)

	if err != nil {

		return nil, err
	}

	defer sql.Close()

	var Data []models.ResponseItem

	for sql.Next() {
		var responseItem models.ResponseItem

		if err := sql.Scan(&responseItem.Id, &responseItem.Title, &responseItem.Code, &responseItem.Value, &responseItem.Adate); err != nil {
			// handle error
			return nil, err
		}
		Data = append(Data, responseItem)
	}

	return Data, nil
}

const getByDateCodeQuery = "SELECT * FROM  r_currency  where code = ? AND A_DATE = str_to_date(?,'%d.%m.%Y')"

func (repo *MySQLUserRepository) GetByDateCode(ctx context.Context, date string, code string) ([]models.ResponseItem, error) {
	// Implementation
	sql, err := repo.Db.QueryContext(ctx, getByDateCodeQuery, code, date)

	if err != nil {

		return nil, err
	}
	defer sql.Close()

	var v struct {
		Data []models.ResponseItem `json:"data"`
	}

	for sql.Next() {
		var responseItem models.ResponseItem

		if err := sql.Scan(&responseItem.Id, &responseItem.Title, &responseItem.Code, &responseItem.Value, &responseItem.Adate); err != nil {
			// handle error
			return nil, err
		}
		v.Data = append(v.Data, responseItem)
	}

	return v.Data, nil
}

const existsQuery = "SELECT COUNT(*) FROM  r_currency  where code = ? AND A_DATE = str_to_date(?,'%d.%m.%Y')"

func (repo *MySQLUserRepository) Exists(ctx context.Context, user *models.Item) (int, error) {
	// Implementation
	var count int
	err := repo.Db.QueryRowContext(ctx, existsQuery, user.Title, user.Date).Scan(&count)

	if err != nil {
		return 0, err

	}

	return count, nil
}

const updateQuery = "UPDATE r_currency SET VALUE = ? WHERE  CODE = ? AND  A_DATE = str_to_date(?,'%d.%m.%Y')"

func (repo *MySQLUserRepository) Update(ctx context.Context, user *models.Item) error {
	// Implementation
	_, err := repo.Db.ExecContext(ctx, updateQuery, user.Description, user.Title, user.Date)
	if err != nil {
		return err
	}

	return nil
}

const insertQuery = "INSERT INTO r_currency (TITLE, CODE, VALUE, A_DATE) VALUES (?, ?, ?, str_to_date(?,'%d.%m.%Y'))"

func (repo *MySQLUserRepository) Insert(ctx context.Context, user *models.Item) error {
	// Implementation
	_, err := repo.Db.ExecContext(ctx, insertQuery, user.Fullname, user.Title, user.Description, user.Date)

	if err != nil {
		return err
	}
	return nil
}
