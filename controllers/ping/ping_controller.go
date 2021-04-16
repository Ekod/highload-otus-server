package ping

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//Ping для хост платформы, чтобы проверять состояние сервиса
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
