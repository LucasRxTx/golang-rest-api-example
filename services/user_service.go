package services

import (
	"database/sql"
	"rest-api/database"
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
	transaction := func(tx *sql.Tx) ([]domain.UserFriend, error) {
		return userService.UserFreindsRepo.GetAllForUserId(tx, userId)
	}

	fail := func(err error) ([]domain.UserFriend, error) {
		return []domain.UserFriend{}, err
	}

	return database.TransactionWithValue[[]domain.UserFriend]{
		Db:          userService.Db,
		Transaction: transaction,
		OnError:     fail,
	}.Execute()
}

func (userService *UserService) UpdateFriends(id string, friends []string) error {
	transaction := func(tx *sql.Tx) error {
		return userService.UserFreindsRepo.Update(tx, id, friends)
	}

	fail := func(err error) error {
		return err
	}

	return database.Transaction{
		Db:          userService.Db,
		Transaction: transaction,
		OnError:     fail,
	}.Execute()
}

func (userService *UserService) GetAll() ([]domain.User, error) {
	transaction := func(tx *sql.Tx) ([]domain.User, error) {
		return userService.UserRepo.GetAll(tx)
	}

	fail := func(err error) ([]domain.User, error) {
		return []domain.User{}, err
	}

	return database.TransactionWithValue[[]domain.User]{
		Db:          userService.Db,
		Transaction: transaction,
		OnError:     fail,
	}.Execute()
}

func (userService *UserService) GetUserById(id string) (domain.User, error) {
	transaction := func(tx *sql.Tx) (domain.User, error) {
		return userService.UserRepo.GetById(tx, id)
	}

	fail := func(err error) (domain.User, error) {
		return domain.User{}, err
	}

	return database.TransactionWithValue[domain.User]{
		Db:          userService.Db,
		Transaction: transaction,
		OnError:     fail,
	}.Execute()
}

func (userService *UserService) CreateUser(name string) (uuid.UUID, error) {
	transaction := func(tx *sql.Tx) (uuid.UUID, error) {
		userId, err := userService.UserRepo.Create(tx, name)
		if err != nil {
			return uuid.UUID{}, err
		}

		_, err = userService.GameStateRepo.Create(tx, userId.String(), 0, 0)
		if err != nil {
			return uuid.UUID{}, err
		}

		return userId, nil
	}

	fail := func(err error) (uuid.UUID, error) {
		return uuid.UUID{}, err
	}

	return database.TransactionWithValue[uuid.UUID]{
		Db:          userService.Db,
		Transaction: transaction,
		OnError:     fail,
	}.Execute()
}

func (userService *UserService) GetGameState(userId string) (domain.GameState, error) {
	transaction := func(tx *sql.Tx) (domain.GameState, error) {
		return userService.GameStateRepo.GetByUserId(tx, userId)
	}

	fail := func(err error) (domain.GameState, error) {
		return domain.GameState{}, err
	}

	return database.TransactionWithValue[domain.GameState]{
		Db:          userService.Db,
		Transaction: transaction,
		OnError:     fail,
	}.Execute()
}

func (userService *UserService) UpdateGameState(id string, gamesPlayed int, score int) error {
	transaction := func(tx *sql.Tx) error {
		gameState, err := userService.GameStateRepo.GetByUserId(tx, id)
		if err != nil {
			return err
		}

		return userService.GameStateRepo.Update(tx, gameState.Id.String(), gamesPlayed, score)
	}

	fail := func(err error) error {
		return err
	}

	return database.Transaction{
		Db:          userService.Db,
		Transaction: transaction,
		OnError:     fail,
	}.Execute()
}
