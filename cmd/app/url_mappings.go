package main

import (
	"github.com/Ekod/highload-otus/internal/transport/controllers"
	"github.com/Ekod/highload-otus/internal/transport/controllers/ping"
	"github.com/Ekod/highload-otus/internal/transport/services"
	"github.com/Ekod/highload-otus/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"os"
)

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
	services *services.Services
}

func APIMux(cfg APIMuxConfig) *gin.Engine {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	config.AllowMethods = []string{http.MethodGet, http.MethodPost, http.MethodDelete}
	router.Use(cors.New(config))
	router.Use(gin.Recovery())
	router.GET("/ping", ping.Ping)
	userHandlers := controllers.UserHandlers{
		Service: cfg.services,
	}
	friendHandlers := controllers.FriendHandlers{
		Service: cfg.services,
	}
	apiGroup := router.Group("/api")
	{
		apiGroup.POST("/login", userHandlers.LoginUser)
		apiGroup.POST("/register", userHandlers.RegisterUser)
		apiGroup.GET("/users", middlewares.CheckToken, userHandlers.GetUsers)
		apiGroup.GET("/info", middlewares.CheckToken, userHandlers.GetCurrentUser)
		apiGroup.GET("/search-users", middlewares.CheckToken, userHandlers.GetUsersByFullName)
		apiGroup.GET("/friends", middlewares.CheckToken, friendHandlers.GetFriends)
		apiGroup.POST("/make-friends", middlewares.CheckToken, friendHandlers.MakeFriends)
		apiGroup.DELETE("/remove-friend_service/:id", middlewares.CheckToken, friendHandlers.RemoveFriend)
	}

	return router
}
