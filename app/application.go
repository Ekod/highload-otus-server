package app

import (
	_ "github.com/Ekod/highload-otus/datasources/mysql/users_db"
	_ "github.com/Ekod/highload-otus/utils/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

var router = gin.Default()

//StartApplication задаёт изначальный конфиг CORS и запускает сервис
func StartApplication() {
	port := os.Getenv("PORT")
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	config.AllowMethods = []string{http.MethodGet, http.MethodPost, http.MethodDelete}
	router.Use(cors.New(config))
	router.Use(gin.Recovery())
	mapUrls()
	if err := router.Run(":" + port); err != nil {
		panic(err)
	}
}
