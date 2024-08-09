package db

import (
	"currency/internal/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func NewDb(aconfig *config.Config) *sql.DB {

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
