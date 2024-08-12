package main

import (
	"fmt"
	"rest-api/controllers"
	"rest-api/database"
	"rest-api/repository"
	"rest-api/services"
	"rest-api/settings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {

	settings.Validate()

	db, err := database.GetConnection()
	if err != nil {
		panic(fmt.Sprintf("Could not get database connection: %s\n", err.Error()))
	}

	userService := services.UserService{ // make constructor
		Db:              db,
		UserRepo:        &repository.UserRepository{},
		GameStateRepo:   &repository.GameStateRepository{},
		UserFreindsRepo: &repository.UserFriendsRepository{},
	}

	userController := controllers.UserController{
		UserService: &userService,
	}

	gameStateController := controllers.GameStateConroller{
		UserService: &userService,
	}

	userFriendsController := controllers.UserFriendsController{
		UserService: &userService,
	}

	router := gin.Default()
	router.POST("/user/", userController.CreateUser)
	router.GET("/user/", userController.GetUsers)
	router.GET("/user/:id", userController.GetUserById)
	router.GET("/user/:id/state", gameStateController.GetSavedGameByUserId)
	router.POST("/user/:id/state", gameStateController.SaveGameState)
	router.GET("/user/:id/friends", userFriendsController.GetFriendsByUserId)
	router.POST("/user/:id/friends", userFriendsController.UpdateFriends)

	router.Run("0.0.0.0:8080")
}
