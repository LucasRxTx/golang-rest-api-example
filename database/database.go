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

type Transaction struct {
	Db          *sql.DB
	Transaction func(tx *sql.Tx) error
	OnError     func(err error) error
}

func (t Transaction) Execute() error {
	tx, err := t.Db.Begin()
	if err != nil {
		return t.OnError(err)
	}

	defer tx.Rollback()

	err = t.Transaction(tx)
	if err != nil {
		return t.OnError(err)
	}

	err = tx.Commit()

	return err
}

type TransactionWithValue[T any] struct {
	Db          *sql.DB
	Transaction func(tx *sql.Tx) (T, error)
	OnError     func(err error) (T, error)
}

func (t TransactionWithValue[T]) Execute() (T, error) {
	tx, err := t.Db.Begin()
	if err != nil {
		return t.OnError(err)
	}

	defer tx.Rollback()

	value, err := t.Transaction(tx)
	if err != nil {
		return t.OnError(err)
	}

	err = tx.Commit()
	if err != nil {
		return t.OnError(err)
	}

	return value, nil
}
