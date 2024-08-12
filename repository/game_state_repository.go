package repository

import (
	"database/sql"
	"errors"
	"rest-api/domain"

	"github.com/google/uuid"
)

type IGameStateRepository interface {
	Create(tx *sql.Tx, userId string, gamesPlayed int, score int) (uuid.UUID, error)
	Update(tx *sql.Tx, gameId string, gamesPlayed int, score int) error
	GetById(tx *sql.Tx, id string) (domain.GameState, error)
	GetByUserId(tx *sql.Tx, id string) (domain.GameState, error)
}

type GameStateRepository struct{}

func (repo *GameStateRepository) GetById(tx *sql.Tx, id string) (domain.GameState, error) {
	rows, err := tx.Query(
		`SELECT * FROM game_state WHERE id = $1 LIMIT 1;`,
		id,
	)

	if err != nil {
		return domain.GameState{}, err
	}

	defer rows.Close()

	var gameState domain.GameState
	for rows.Next() {
		err = rows.Scan(&gameState.Id, &gameState.UserId, &gameState.GamesPlayed, &gameState.Score)
		if err != nil {
			return domain.GameState{}, err
		}
	}

	err = rows.Err()
	if err != nil {
		return domain.GameState{}, err
	}

	if gameState.Id.String() == "" {
		return domain.GameState{}, errors.New("user not found")
	}

	return gameState, nil
}

func (repo *GameStateRepository) GetByUserId(tx *sql.Tx, id string) (domain.GameState, error) {
	rows, err := tx.Query(
		`SELECT * FROM game_state WHERE user_id = $1 LIMIT 1;`,
		id,
	)

	if err != nil {
		return domain.GameState{}, err
	}

	defer rows.Close()

	var gameState domain.GameState
	for rows.Next() {
		err = rows.Scan(&gameState.Id, &gameState.UserId, &gameState.GamesPlayed, &gameState.Score)
		if err != nil {
			return domain.GameState{}, err
		}
	}

	err = rows.Err()
	if err != nil {
		return domain.GameState{}, err
	}

	if gameState.Id.String() == "" {
		return domain.GameState{}, errors.New("user not found")
	}

	return gameState, nil
}

func (repo *GameStateRepository) Update(tx *sql.Tx, gameId string, gamesPlayed int, score int) error {
	_, err := tx.Exec(
		`UPDATE game_state
		SET
			("games_played", "score") = ($1, $2)
		WHERE
			id = $3
		;`,
		gamesPlayed,
		score,
		gameId,
	)

	return err
}

func (repo *GameStateRepository) Create(tx *sql.Tx, userId string, gamesPlayed int, score int) (uuid.UUID, error) {
	gameId := uuid.New()
	_, err := tx.Exec(
		`INSERT INTO game_state ("id", "user_id", "games_played", "score") values ($1, $2, $3, $4);`,
		gameId,
		userId,
		gamesPlayed,
		score,
	)

	if err != nil {
		return uuid.UUID{}, err
	}

	return gameId, nil
}
