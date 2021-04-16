package controllers

import (
	"net/http"
	"strconv"

	"github.com/Ekod/highload-otus/domain/users"
	"github.com/Ekod/highload-otus/services"
	"github.com/Ekod/highload-otus/utils/errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var UserController userContrlInterface = &userController{}

type userController struct{}

type userContrlInterface interface {
	LoginUser(c *gin.Context)
	RegisterUser(c *gin.Context)
	GetUsers(c *gin.Context)
	GetCurrentUser(c *gin.Context)
	GetFriends(c *gin.Context)
	MakeFriends(c *gin.Context)
	RemoveFriend(c *gin.Context)
}

func (uc *userController) RemoveFriend(c *gin.Context) {
	userId, err := services.SecurityService.GetUserIdFromToken(c)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	friendId := c.Param("id")
	var err2 error
	parsedFriendId, err2 := strconv.Atoi(friendId)
	if err2 != nil {
		log.Error("Controllers_RemoveFriend - Error parsing incoming JSON")
		restErr := errors.NewBadRequestError("Invalid params")
		c.JSON(restErr.Status, restErr)
		return
	}
	response, err := services.UserService.RemoveFriend(userId, int64(parsedFriendId))
	c.JSON(http.StatusOK, response)
}

//LoginUser логинит пользователя
func (uc *userController) LoginUser(c *gin.Context) {
	var user users.UserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Error("Controllers_LoginUser - Error parsing incoming JSON")
		restErr := errors.NewBadRequestError("Invalid json request")
		c.JSON(restErr.Status, restErr)
		return
	}

	response, loginErr := services.UserService.LoginUser(&user)
	if loginErr != nil {
		c.JSON(loginErr.Status, loginErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

//RegisterUser регистрирует пользователя
func (uc *userController) RegisterUser(c *gin.Context) {
	var user users.UserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Error("Controllers_RegisterUser - Error parsing incoming JSON")
		restErr := errors.NewBadRequestError("Invalid json request")
		c.JSON(restErr.Status, restErr)
		return
	}

	response, saveErr := services.UserService.RegisterUser(&user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		return
	}

	c.JSON(http.StatusCreated, response)
}

//GetUsers для получения всех пользователей
func (uc *userController) GetUsers(c *gin.Context) {
	response, err := services.UserService.GetUsers()
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

//GetCurrentUser для получения инфы по пользователю
func (uc *userController) GetCurrentUser(c *gin.Context) {
	userId, err := services.SecurityService.GetUserIdFromToken(c)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	response, err := services.UserService.GetCurrentUser(userId)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (uc *userController) GetFriends(c *gin.Context) {
	userId, err := services.SecurityService.GetUserIdFromToken(c)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	response, err := services.UserService.GetFriends(userId)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (uc *userController) MakeFriends(c *gin.Context) {
	var friend users.UserFriend
	if err := c.ShouldBindJSON(&friend); err != nil {
		log.Error("Controllers_RegisterUser - Error parsing incoming JSON")
		restErr := errors.NewBadRequestError("Invalid json request")
		c.JSON(restErr.Status, restErr)
		return
	}
	userId, err := services.SecurityService.GetUserIdFromToken(c)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	response, err := services.UserService.MakeFriends(userId, &friend)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
