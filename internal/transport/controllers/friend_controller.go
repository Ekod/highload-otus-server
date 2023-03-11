package controllers

import (
	"github.com/Ekod/highload-otus/internal/transport/services"
	"github.com/Ekod/highload-otus/utils/errors"
	"github.com/Ekod/highload-otus/utils/security"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type FriendHandlers struct {
	Service *services.Services
	Logger  *zap.SugaredLogger
}

func (h *FriendHandlers) RemoveFriend(c *gin.Context) {
	ctx := c.Request.Context()

	userId, err := security.GetUserIdFromToken(c)
	if err != nil {
		h.Logger.Error(err)

		err := errors.ParseError(err)

		c.JSON(err.Status, err)
		return
	}

	friendId := c.Param("id")

	parsedFriendId, err := strconv.Atoi(friendId)
	if err != nil {
		h.Logger.Error("[ERROR] Controllers_RemoveFriend - Error parsing incoming JSON")

		err := errors.NewHandlerBadRequestError("Invalid params")
		c.JSON(err.Status, err)

		return
	}

	err = h.Service.FriendsService.RemoveFriend(ctx, userId, parsedFriendId)

	c.Status(http.StatusOK)
}

func (h *FriendHandlers) GetFriends(c *gin.Context) {
	ctx := c.Request.Context()

	userId, err := security.GetUserIdFromToken(c)
	if err != nil {
		h.Logger.Error(err)

		err := errors.ParseError(err)

		c.JSON(err.Status, err)
		return
	}

	response, err := h.Service.FriendsService.GetFriends(ctx, userId)
	if err != nil {
		h.Logger.Error(err)

		err := errors.ParseError(err)

		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *FriendHandlers) MakeFriends(c *gin.Context) {
	ctx := c.Request.Context()

	var friendID int
	if err := c.ShouldBindJSON(&friendID); err != nil {
		h.Logger.Error("[ERROR] Controllers_RegisterUser - Error parsing incoming JSON")

		err := errors.NewHandlerBadRequestError("Invalid json request")
		c.JSON(err.Status, err)

		return
	}

	userId, err := security.GetUserIdFromToken(c)
	if err != nil {
		h.Logger.Error(err)

		err := errors.ParseError(err)

		c.JSON(err.Status, err)
		return
	}

	response, err := h.Service.FriendsService.MakeFriends(ctx, userId, friendID)
	if err != nil {
		h.Logger.Error(err)

		err := errors.ParseError(err)

		c.JSON(err.Status, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
