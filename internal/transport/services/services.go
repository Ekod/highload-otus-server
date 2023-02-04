package services

import (
	"github.com/Ekod/highload-otus/domain/users"
	"github.com/Ekod/highload-otus/utils/errors"
)

type Services struct {
	UserService    UserService
	FriendsService FriendsService
}

func New(userService UserService, friendsService FriendsService) *Services {
	return &Services{
		UserService:    userService,
		FriendsService: friendsService,
	}
}

type UserService interface {
	LoginUser(user *users.UserRequest) (map[string]interface{}, *errors.RestErr)
	RegisterUser(user *users.UserRequest) (map[string]interface{}, *errors.RestErr)
	GetUsers() (map[string]interface{}, *errors.RestErr)
	GetCurrentUser(userId int64) (map[string]interface{}, *errors.RestErr)
	GetUsersByFullName(firstName, lastName string) (map[string][]users.ResponseUser, *errors.RestErr)
}

type FriendsService interface {
	GetFriends(userId int64) (map[string]interface{}, *errors.RestErr)
	MakeFriends(userId int64, friend *users.UserFriend) (map[string]interface{}, *errors.RestErr)
	RemoveFriend(userId int64, friendId int64) (map[string]interface{}, *errors.RestErr)
}
