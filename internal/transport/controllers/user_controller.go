package controllers

import (
	"github.com/Ekod/highload-otus/domain/users"
	"github.com/Ekod/highload-otus/internal/transport/services"
	"github.com/Ekod/highload-otus/utils/errors"
	"github.com/Ekod/highload-otus/utils/security"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type UserHandlers struct {
	Service *services.Services
	Logger  *zap.SugaredLogger
}

//RegisterUser регистрирует пользователя
func (h *UserHandlers) RegisterUser(c *gin.Context) {
	ctx := c.Request.Context()

	var user users.UserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		h.Logger.Error("[ERROR] Controllers_RegisterUser - Error parsing incoming JSON")

		err := errors.NewHandlerBadRequestError("Invalid json request")
		c.JSON(err.Status, err)

		return
	}

	response, err := h.Service.UserService.RegisterUser(ctx, &user)
	if err != nil {
		h.Logger.Error(err)

		err := errors.ParseError(err)

		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

//GetUsers для получения всех пользователей
func (h *UserHandlers) GetUsers(c *gin.Context) {
	ctx := c.Request.Context()

	response, err := h.Service.UserService.GetUsers(ctx)
	if err != nil {
		h.Logger.Error(err)

		err := errors.ParseError(err)

		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

//GetCurrentUser для получения инфы по пользователю
func (h *UserHandlers) GetCurrentUser(c *gin.Context) {
	ctx := c.Request.Context()

	userId, err := security.GetUserIdFromToken(c)
	if err != nil {
		h.Logger.Error(err)

		err := errors.ParseError(err)

		c.JSON(err.Status, err)
		return
	}

	response, err := h.Service.UserService.GetCurrentUser(ctx, userId)
	if err != nil {
		h.Logger.Error(err)

		err := errors.ParseError(err)

		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

//LoginUser логинит пользователя
func (h *UserHandlers) LoginUser(c *gin.Context) {
	ctx := c.Request.Context()

	var user users.UserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		h.Logger.Error("[ERROR] Controllers_LoginUser - Error parsing incoming JSON")

		err := errors.NewHandlerBadRequestError("Invalid json request")
		c.JSON(err.Status, err)

		return
	}

	response, err := h.Service.UserService.LoginUser(ctx, &user)
	if err != nil {
		h.Logger.Error(err)

		err := errors.ParseError(err)

		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandlers) GetUsersByFullName(c *gin.Context) {
	ctx := c.Request.Context()

	firstName := c.Query("firstName")
	lastName := c.Query("lastName")

	response, err := h.Service.UserService.GetUsersByFullName(ctx, firstName, lastName)
	if err != nil {
		h.Logger.Error(err)

		err := errors.ParseError(err)

		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
