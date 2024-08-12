package database

import (
	"database/sql"
	"fmt"
	"rest-api/settings"
)

func GetConnection() (*sql.DB, error) {
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", settings.DB_HOST, settings.DB_PORT, settings.DB_USER, settings.DB_PASSWORD, settings.DB_DATABASE)
	return sql.Open("postgres", psqlconn)
}
