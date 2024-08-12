package controllers

import (
	"rest-api/dto"
	"rest-api/services"

	"github.com/gin-gonic/gin"
)

type GameStateConroller struct {
	UserService services.IUserService
}

type CreateGameStateRequest struct {
	GamesPlayed int `json:"gamesPlayed" binding:"required"`
	Score       int `json:"score" binding:"required"`
}

func (controller *GameStateConroller) SaveGameState(c *gin.Context) {
	var uriParams GetUserByIdUriParams
	err := c.ShouldBindUri(&uriParams)
	if err != nil {
		c.JSON(400, dto.ErrorDto{
			Message: "Invalid user id",
		})
		return
	}

	var request CreateGameStateRequest
	err = c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(400, dto.ErrorDto{
			Message: "Invalid game state",
		})
		return
	}

	err = controller.UserService.UpdateGameState(uriParams.Id, request.GamesPlayed, request.Score)
	if err != nil {
		c.JSON(500, dto.ErrorDto{
			Message: err.Error(),
		})
		return
	}

	c.JSON(201, dto.GameStateDto{
		GamesPlayed: request.GamesPlayed,
		Score:       request.Score,
	})
}

func (controller *GameStateConroller) GetSavedGameByUserId(c *gin.Context) {
	var uriParams GetUserByIdUriParams
	err := c.ShouldBindUri(&uriParams)
	if err != nil {
		c.JSON(400, dto.ErrorDto{
			Message: "Invalid user id",
		})
		return
	}

	gameState, err := controller.UserService.GetGameState(uriParams.Id)
	if err != nil {
		c.JSON(500, dto.ErrorDto{
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, dto.GameStateDto{
		GamesPlayed: gameState.GamesPlayed,
		Score:       gameState.Score,
	})
}
