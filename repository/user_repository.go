package repository

import (
	"database/sql"
	"rest-api/domain"

	"github.com/google/uuid"
)

type IUserRepository interface {
	Create(tx *sql.Tx, name string) (uuid.UUID, error)
	GetById(tx *sql.Tx, id string) (domain.User, error)
	GetAll(tx *sql.Tx) ([]domain.User, error)
}

type UserRepository struct{}

func (repo *UserRepository) Create(tx *sql.Tx, userName string) (uuid.UUID, error) {
	var newUserId = uuid.New()
	_, err := tx.Exec(
		`INSERT INTO game_user ("id", "name") values ($1, $2);`,
		newUserId,
		userName,
	)

	if err != nil {
		return uuid.UUID{}, err
	}

	return newUserId, nil
}

func (repo *UserRepository) GetById(tx *sql.Tx, id string) (domain.User, error) {
	rows := tx.QueryRow(`SELECT * FROM game_user WHERE id = $1 LIMIT 1;`, id)

	var user domain.User
	err := rows.Scan(&user.Id, &user.Name)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (repo *UserRepository) GetAll(tx *sql.Tx) ([]domain.User, error) {
	rows, err := tx.Query(`SELECT * FROM game_user;`)
	if err != nil {
		return []domain.User{}, err
	}

	defer rows.Close()

	var users []domain.User
	var user domain.User
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Name)
		if err != nil {
			return []domain.User{}, err
		}
		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		return []domain.User{}, err
	}

	return users, nil
}
