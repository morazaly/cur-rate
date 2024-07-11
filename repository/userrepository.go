package repository

import (
	"currency/models"
	"database/sql"
)

type MySQLUserRepository struct {
	// DB connection and other fields
	Db *sql.DB
}

func (repo *MySQLUserRepository) GetByDate(Date string) ([]models.ResponseItem, error) {
	// Implementation

	sql, err := repo.Db.Query("SELECT * FROM  r_currency  where  A_DATE = str_to_date(?,'%d.%m.%Y')", Date)

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

func (repo *MySQLUserRepository) GetByDateCode(Date string, Code string) ([]models.ResponseItem, error) {
	// Implementation
	sql, err := repo.Db.Query("SELECT * FROM  r_currency  where code = ? AND A_DATE = str_to_date(?,'%d.%m.%Y')", Code, Date)

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
func (repo *MySQLUserRepository) Exists(user *models.Item) (int, error) {
	// Implementation
	var count int
	err := repo.Db.QueryRow("SELECT COUNT(*) FROM  r_currency  where code = ? AND A_DATE = str_to_date(?,'%d.%m.%Y')", user.Title, user.Date).Scan(&count)

	if err != nil {
		return 0, err

	}

	return count, nil
}
func (repo *MySQLUserRepository) Update(user *models.Item) error {
	// Implementation
	_, err := repo.Db.Exec("UPDATE r_currency SET VALUE = ? WHERE  CODE = ? AND  A_DATE = str_to_date(?,'%d.%m.%Y')", user.Description, user.Title, user.Date)
	if err != nil {
		return err
	}

	return nil
}
func (repo *MySQLUserRepository) Insert(user *models.Item) error {
	// Implementation
	_, err := repo.Db.Exec("INSERT INTO r_currency (TITLE, CODE, VALUE, A_DATE) VALUES (?, ?, ?, str_to_date(?,'%d.%m.%Y'))", user.Fullname, user.Title, user.Description, user.Date)

	if err != nil {
		return err
	}
	return nil
}
