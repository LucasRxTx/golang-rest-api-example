package settings

import (
	"fmt"
	"os"
)

var (
	APP_NAME    = os.Getenv("APP_NAME")
	DB_DATABASE = os.Getenv("DB_DATABASE")
	DB_HOST     = os.Getenv("DB_HOST")
	DB_PORT     = os.Getenv("DB_PORT")
	DB_USER     = os.Getenv("DB_USER")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
)

func Validate() error {
	var api_settings = map[string]string{
		"APP_NAME":    APP_NAME,
		"DB_DATABASE": DB_DATABASE,
		"DB_HOST":     DB_HOST,
		"DB_PORT":     DB_PORT,
		"DB_USER":     DB_USER,
		"DB_PASSWORD": DB_PASSWORD,
	}

	for key, value := range api_settings {
		if value == "" {
			return fmt.Errorf("settings: %s is not set in the environment and is required", key)
		}
	}

	return nil
}
