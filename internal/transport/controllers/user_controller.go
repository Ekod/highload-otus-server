package controllers

import (
	"github.com/Ekod/highload-otus/domain/users"
	"github.com/Ekod/highload-otus/internal/transport/services"
	"github.com/Ekod/highload-otus/utils/errors"
	"github.com/Ekod/highload-otus/utils/security"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type UserHandlers struct {
	Service *services.Services
}

//RegisterUser регистрирует пользователя
func (h *UserHandlers) RegisterUser(c *gin.Context) {
	var user users.UserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Error("Controllers_RegisterUser - Error parsing incoming JSON")
		restErr := errors.NewBadRequestError("Invalid json request")
		c.JSON(restErr.Status, restErr)
		return
	}

	response, saveErr := h.Service.UserService.RegisterUser(&user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}

	c.JSON(http.StatusCreated, response)
}

//GetUsers для получения всех пользователей
func (h *UserHandlers) GetUsers(c *gin.Context) {
	response, err := h.Service.UserService.GetUsers()
	if err != nil {
		log.Error(err)
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

//GetCurrentUser для получения инфы по пользователю
func (h *UserHandlers) GetCurrentUser(c *gin.Context) {
	userId, err := security.GetUserIdFromToken(c)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	response, err := h.Service.UserService.GetCurrentUser(userId)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

//LoginUser логинит пользователя
func (h *UserHandlers) LoginUser(c *gin.Context) {
	var user users.UserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Error("Controllers_LoginUser - Error parsing incoming JSON")
		restErr := errors.NewBadRequestError("Invalid json request")
		c.JSON(restErr.Status, restErr)
		return
	}

	response, loginErr := h.Service.UserService.LoginUser(&user)
	if loginErr != nil {
		c.JSON(loginErr.Status, loginErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandlers) GetUsersByFullName(c *gin.Context) {
	firstName := c.Query("firstName")
	lastName := c.Query("lastName")
	response, err := h.Service.UserService.GetUsersByFullName(firstName, lastName)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, response)
}
