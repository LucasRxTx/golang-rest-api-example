package main

import (
	"log"
	"rest-api/controllers"
	"rest-api/database"
	"rest-api/repository"
	"rest-api/services"
	"rest-api/settings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	err := settings.Validate()
	if err != nil {
		log.Fatal("Invalid settings for app:", err.Error())
	}

	db, err := database.GetConnection()
	if err != nil {
		log.Fatal("Could not connect to database: ", err)
	}

	userService := services.UserService{
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
