package db

import (
	"currency/internal/models"
	"database/sql"
	"fmt"
	"log"
)

func NewDb(aconfig *models.Config) *sql.DB {

	// Подключение к базе данных MySQL
	dataSourceName := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		aconfig.MySQL.User,
		aconfig.MySQL.Password,
		aconfig.MySQL.Host,
		aconfig.MySQL.Port,
		aconfig.MySQL.Database,
	)

	// Подключение к базе данных MySQL
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	// Проверка подключения к базе данных
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}
