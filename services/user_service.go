package services

import (
	"database/sql"
	"rest-api/dao"
	"rest-api/repository"

	"github.com/google/uuid"
)

type IUserService interface {
	CreateUser(name string) (uuid.UUID, error)
	GetUserById(id string) (dao.UserDao, error)
	GetAll() ([]dao.UserDao, error)
	GetGameState(id string) (dao.GameStateDao, error)
	UpdateGameState(id string, gamesPlayer int, score int) error
	UpdateFriends(id string, friends []string) error
	GetAllFriends(id string) ([]dao.UserFriendDao, error)
}

type UserService struct {
	Db              *sql.DB
	UserRepo        repository.IUserRepository
	GameStateRepo   repository.IGameStateRepository
	UserFreindsRepo repository.IUserFriendsRepo
}

func (userService *UserService) GetAllFriends(userId string) ([]dao.UserFriendDao, error) {
	fail := func(err error) ([]dao.UserFriendDao, error) {
		return []dao.UserFriendDao{}, err
	}

	tx, err := userService.Db.Begin()
	if err != nil {
		return fail(err)
	}

	defer tx.Rollback()

	userFriends, err := userService.UserFreindsRepo.GetAllForUserId(tx, userId)
	if err != nil {
		return fail(err)
	}

	err = tx.Commit()
	if err != nil {
		return fail(err)
	}

	return userFriends, nil
}

func (userService *UserService) UpdateFriends(id string, friends []string) error {
	tx, err := userService.Db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()
	err = userService.UserFreindsRepo.Update(tx, id, friends)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (userService *UserService) GetAll() ([]dao.UserDao, error) {
	fail := func(err error) ([]dao.UserDao, error) {
		return []dao.UserDao{}, err
	}

	tx, err := userService.Db.Begin()
	if err != nil {
		return fail(err)
	}

	defer tx.Rollback()

	users, err := userService.UserRepo.GetAll(tx)
	if err != nil {
		return fail(err)
	}

	err = tx.Commit()
	if err != nil {
		return fail(err)
	}

	return users, nil
}

func (userService *UserService) GetUserById(id string) (dao.UserDao, error) {
	fail := func(err error) (dao.UserDao, error) {
		return dao.UserDao{}, err
	}

	tx, err := userService.Db.Begin()
	if err != nil {
		return fail(err)
	}

	defer tx.Rollback()

	user, err := userService.UserRepo.GetById(tx, id)
	if err != nil {
		return fail(err)
	}

	err = tx.Commit()
	if err != nil {
		return fail(err)
	}

	return user, nil
}

func (userService *UserService) CreateUser(name string) (uuid.UUID, error) {
	fail := func(err error) (uuid.UUID, error) {
		return uuid.UUID{}, err
	}

	tx, err := userService.Db.Begin()
	if err != nil {
		return fail(err)
	}

	defer tx.Rollback()

	uuid, err := userService.UserRepo.Create(tx, name)
	if err != nil {
		return fail(err)
	}

	_, err = userService.GameStateRepo.Create(tx, uuid.String(), 0, 0)
	if err != nil {
		return fail(err)
	}

	err = tx.Commit()
	if err != nil {
		return fail(err)
	}

	return uuid, nil
}

func (userService *UserService) GetGameState(userId string) (dao.GameStateDao, error) {
	fail := func(err error) (dao.GameStateDao, error) {
		return dao.GameStateDao{}, err
	}

	tx, err := userService.Db.Begin()
	if err != nil {
		return fail(err)
	}

	defer tx.Rollback()

	gameState, err := userService.GameStateRepo.GetByUserId(tx, userId)
	if err != nil {
		return fail(err)
	}

	err = tx.Commit()
	if err != nil {
		return fail(err)
	}

	return gameState, nil
}

func (userService *UserService) UpdateGameState(id string, gamesPlayed int, score int) error {
	tx, err := userService.Db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	gameState, err := userService.GameStateRepo.GetByUserId(tx, id)
	if err != nil {
		return err
	}

	err = userService.GameStateRepo.Update(tx, gameState.Id.String(), gamesPlayed, score)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
