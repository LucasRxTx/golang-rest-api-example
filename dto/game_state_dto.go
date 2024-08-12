package dto

type GameStateDto struct {
	GamesPlayed int `json:"gamesPlayed" binding:"required"`
	Score       int `json:"score" binding:"required"`
}
