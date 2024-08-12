package controllers

import (
	"rest-api/dto"
	"rest-api/services"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type UserController struct {
	UserService services.IUserService
}

type CreateUserRequest struct {
	Name string `json:"name" binding:"required"`
}

func (controller *UserController) CreateUser(c *gin.Context) {
	var createUserRequest CreateUserRequest
	err := c.ShouldBindJSON(&createUserRequest)
	if err != nil {
		c.JSON(400, dto.ErrorDto{
			Message: err.Error(),
		})
		return
	}

	userId, err := controller.UserService.CreateUser(createUserRequest.Name)
	if err != nil {
		c.JSON(500, dto.ErrorDto{
			Message: err.Error(),
		})
		return
	}

	user, err := controller.UserService.GetUserById(userId.String())
	if err != nil {
		c.JSON(500, dto.ErrorDto{
			Message: err.Error(),
		})
		return
	}

	c.JSON(201, dto.UserDto{
		Id:   user.Id.String(),
		Name: user.Name,
	})
}

func (controller *UserController) GetUsers(c *gin.Context) {
	users, err := controller.UserService.GetAll()
	if err != nil {
		c.JSON(500, dto.ErrorDto{
			Message: err.Error(),
		})
		return
	}

	usersOut := []dto.UserDto{}
	for _, user := range users {
		usersOut = append(usersOut, dto.UserDto{
			Id:   user.Id.String(),
			Name: user.Name,
		})
	}

	c.JSON(200, usersOut)
}

type GetUserByIdUriParams struct {
	Id string `uri:"id" binding:"required,uuid"`
}

func (controller *UserController) GetUserById(c *gin.Context) {
	var uriParams GetUserByIdUriParams
	err := c.ShouldBindUri(&uriParams)
	if err != nil {
		c.JSON(400, dto.ErrorDto{
			Message: "Invalid user id",
		})
		return
	}

	user, err := controller.UserService.GetUserById(uriParams.Id)
	if err != err {
		c.JSON(500, dto.ErrorDto{
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, dto.UserDto{
		Id:   user.Id.String(),
		Name: user.Name,
	})
}
