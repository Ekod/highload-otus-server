package middlewares

import (
	"github.com/Ekod/highload-otus/utils/errors"
	"github.com/Ekod/highload-otus/utils/security"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

//CheckToken проверяет наличие id пользователя в токене
func (m *Middleware) CheckToken(c *gin.Context) {
	authorizationHeader := "Authorization"

	header := c.GetHeader(authorizationHeader)
	if header == "" {
		m.Logger.Error("empty auth header")
		c.AbortWithStatusJSON(http.StatusUnauthorized, errors.NewUnauthorizedError("not authorized"))
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		m.Logger.Error("invalid auth header")
		c.AbortWithStatusJSON(http.StatusUnauthorized, errors.NewUnauthorizedError("not authorized"))
	}

	if len(headerParts[1]) == 0 {
		m.Logger.Error("token is empty")
		c.AbortWithStatusJSON(http.StatusUnauthorized, errors.NewUnauthorizedError("not authorized"))
	}

	userId, err := security.ParseToken(headerParts[1])
	if err != nil {
		m.Logger.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusUnauthorized, "not authorized")
	}

	c.Set("userId", userId)
}
