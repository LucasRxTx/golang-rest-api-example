package repository

import (
	"database/sql"
	"errors"
	"rest-api/dao"

	"github.com/google/uuid"
)

type IUserRepository interface {
	Create(tx *sql.Tx, name string) (uuid.UUID, error)
	GetById(tx *sql.Tx, id string) (dao.UserDao, error)
	GetAll(tx *sql.Tx) ([]dao.UserDao, error)
}

type UserRepository struct {
}

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

func (repo *UserRepository) GetById(tx *sql.Tx, id string) (dao.UserDao, error) {
	rows, err := tx.Query(`SELECT * FROM game_user WHERE id = $1 LIMIT 1;`, id)
	if err != nil {
		return dao.UserDao{}, err
	}

	defer rows.Close()

	var user dao.UserDao
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Name)
		if err != nil {
			return dao.UserDao{}, err
		}
	}

	if user.Id.String() == "" {
		return dao.UserDao{}, errors.New("user_repository: User not found")
	}

	return user, nil
}

func (repo *UserRepository) GetAll(tx *sql.Tx) ([]dao.UserDao, error) {
	rows, err := tx.Query(`SELECT * FROM game_user;`)
	if err != nil {
		return []dao.UserDao{}, err
	}

	defer rows.Close()

	var users []dao.UserDao
	var user dao.UserDao
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Name)
		if err != nil {
			return []dao.UserDao{}, err
		}
		users = append(users, user)
	}

	return users, nil
}
