package middlewares

import (
	"github.com/Ekod/highload-otus/services"
	"github.com/Ekod/highload-otus/utils/errors"
	"github.com/Ekod/highload-otus/utils/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
)

//CheckToken проверяет наличие id пользователя в токене
func CheckToken(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		logger.LogErrorMessage("empty auth header")
		c.AbortWithStatusJSON(http.StatusUnauthorized, errors.NewUnauthorizedError("not authorized"))
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		logger.LogErrorMessage("invalid auth header")
		c.AbortWithStatusJSON(http.StatusUnauthorized, errors.NewUnauthorizedError("not authorized"))
	}

	if len(headerParts[1]) == 0 {
		logger.LogErrorMessage("token is empty")
		c.AbortWithStatusJSON(http.StatusUnauthorized, errors.NewUnauthorizedError("not authorized"))
	}

	userId, err := services.SecurityService.ParseToken(headerParts[1])
	if err != nil {
		logger.LogErrorMessage(err.Message)
		c.AbortWithStatusJSON(http.StatusUnauthorized, "not authorized")
	}

	c.Set("userId", userId)
}
