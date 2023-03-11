package services

import (
	"context"
	"github.com/Ekod/highload-otus/domain/users"
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
	LoginUser(ctx context.Context, user *users.UserRequest) (*users.UserResponse, error)
	RegisterUser(ctx context.Context, user *users.UserRequest) (*users.UserResponse, error)
	GetUsers(ctx context.Context) ([]users.UserResponse, error)
	GetCurrentUser(ctx context.Context, userId int) (*users.UserResponse, error)
	GetUsersByFullName(ctx context.Context, firstName, lastName string) ([]users.UserResponse, error)
}

type FriendsService interface {
	GetFriends(ctx context.Context, userId int) ([]users.UserFriend, error)
	MakeFriends(ctx context.Context, userId int, friendID int) (int, error)
	RemoveFriend(ctx context.Context, userId int, friendId int) error
}
