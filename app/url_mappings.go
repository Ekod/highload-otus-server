package app

import (
	"github.com/Ekod/highload-otus/controllers"
	"github.com/Ekod/highload-otus/controllers/ping"
	"github.com/Ekod/highload-otus/middlewares"
)

func mapUrls() {
	router.GET("/ping", ping.Ping)
	apiGroup := router.Group("/api")
	{
		apiGroup.POST("/login", controllers.UserController.LoginUser)
		apiGroup.POST("/register", controllers.UserController.RegisterUser)
		apiGroup.GET("/users", middlewares.CheckToken, controllers.UserController.GetUsers)
		apiGroup.GET("/info", middlewares.CheckToken, controllers.UserController.GetCurrentUser)
		apiGroup.GET("/friends", middlewares.CheckToken, controllers.UserController.GetFriends)
		apiGroup.POST("/make-friends", middlewares.CheckToken, controllers.UserController.MakeFriends)
		apiGroup.DELETE("/remove-friend/:id", middlewares.CheckToken, controllers.UserController.RemoveFriend)

	}

}
