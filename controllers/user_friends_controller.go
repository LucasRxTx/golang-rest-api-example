package controllers

import (
	"rest-api/dto"
	"rest-api/services"

	"github.com/gin-gonic/gin"
)

type UserFriendsController struct {
	UserService services.IUserService
}

func (controller *UserFriendsController) GetFriendsByUserId(c *gin.Context) {
	var uriParams GetUserByIdUriParams
	err := c.ShouldBindUri(&uriParams)
	if err != nil {
		c.JSON(400, dto.ErrorDto{
			Message: "Invalid user id",
		})
		return
	}

	userFriends, err := controller.UserService.GetAllFriends(uriParams.Id)

	if err != nil {
		c.JSON(500, dto.ErrorDto{
			Message: err.Error(),
		})
		return
	}

	friendsList := dto.UserFriendsListDto{}
	friendsList.Freinds = []dto.UserFriendDto{} // todo: This shouldn't be manual
	for _, friend := range userFriends {
		friendsList.AddFriend(dto.UserFriendDto{
			Id:        friend.Id.String(),
			Name:      friend.Name,
			Highscore: friend.Highscore,
		})
	}

	c.JSON(200, friendsList)
}

type UpdateFriendsRequest struct {
	Friends []string `json:"friends" binding:"required"`
}

func (controller *UserFriendsController) UpdateFriends(c *gin.Context) {
	var uriParams GetUserByIdUriParams
	err := c.ShouldBindUri(&uriParams)
	if err != nil {
		c.JSON(400, dto.ErrorDto{
			Message: "Invalid user id",
		})
		return
	}

	var request UpdateFriendsRequest
	err = c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(400, dto.ErrorDto{
			Message: err.Error(),
		})
		return
	}
	err = controller.UserService.UpdateFriends(uriParams.Id, request.Friends)

	if err != nil {
		c.JSON(500, dto.ErrorDto{
			Message: err.Error(),
		})
		return
	}

	c.JSON(201, gin.H{})
}
