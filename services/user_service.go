package services

import (
	"database/sql"
	"rest-api/domain"
	"rest-api/repository"

	"github.com/google/uuid"
)

type IUserService interface {
	CreateUser(name string) (uuid.UUID, error)
	GetUserById(id string) (domain.User, error)
	GetAll() ([]domain.User, error)
	GetGameState(id string) (domain.GameState, error)
	UpdateGameState(id string, gamesPlayer int, score int) error
	UpdateFriends(id string, friends []string) error
	GetAllFriends(id string) ([]domain.UserFriend, error)
}

type UserService struct {
	Db              *sql.DB
	UserRepo        repository.IUserRepository
	GameStateRepo   repository.IGameStateRepository
	UserFreindsRepo repository.IUserFriendsRepo
}

func (userService *UserService) GetAllFriends(userId string) ([]domain.UserFriend, error) {
	fail := func(err error) ([]domain.UserFriend, error) {
		return []domain.UserFriend{}, err
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

func (userService *UserService) GetAll() ([]domain.User, error) {
	fail := func(err error) ([]domain.User, error) {
		return []domain.User{}, err
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

func (userService *UserService) GetUserById(id string) (domain.User, error) {
	fail := func(err error) (domain.User, error) {
		return domain.User{}, err
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

func (userService *UserService) GetGameState(userId string) (domain.GameState, error) {
	fail := func(err error) (domain.GameState, error) {
		return domain.GameState{}, err
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
