package domain

import "github.com/google/uuid"

type GameState struct {
	Id          uuid.UUID `json:"id" binding:"required,uuid"`
	UserId      uuid.UUID `json:"userId" binding:"required,uuid"`
	GamesPlayed int       `json:"gamesPlayed" binding:"required"`
	Score       int       `json:"score" binding:"required"`
}
