package controllers

import (
	"github.com/Ekod/highload-otus/domain/users"
	"github.com/Ekod/highload-otus/internal/transport/services"
	"github.com/Ekod/highload-otus/utils/errors"
	"github.com/Ekod/highload-otus/utils/security"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type FriendHandlers struct {
	Service *services.Services
}

func (h *FriendHandlers) RemoveFriend(c *gin.Context) {
	userId, err := security.GetUserIdFromToken(c)
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
	convertedFriendId := int64(parsedFriendId)

	response, err := h.Service.FriendsService.RemoveFriend(userId, convertedFriendId)
	c.JSON(http.StatusOK, response)
}

func (h *FriendHandlers) GetFriends(c *gin.Context) {
	userId, err := security.GetUserIdFromToken(c)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	response, err := h.Service.FriendsService.GetFriends(userId)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *FriendHandlers) MakeFriends(c *gin.Context) {
	var friend users.UserFriend
	if err := c.ShouldBindJSON(&friend); err != nil {
		log.Error("Controllers_RegisterUser - Error parsing incoming JSON")
		restErr := errors.NewBadRequestError("Invalid json request")
		c.JSON(restErr.Status, restErr)
		return
	}
	userId, err := security.GetUserIdFromToken(c)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	response, err := h.Service.FriendsService.MakeFriends(userId, &friend)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
