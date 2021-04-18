package services

import (
	"github.com/Ekod/highload-otus/domain/users"
	"github.com/Ekod/highload-otus/utils/errors"
	"github.com/Ekod/highload-otus/utils/response"
)

var UserService userServiceInterface = &userService{}

type userService struct{}

type userServiceInterface interface {
	LoginUser(*users.UserRequest) (map[string]interface{}, *errors.RestErr)
	RegisterUser(*users.UserRequest) (map[string]interface{}, *errors.RestErr)
	GetUsers() (map[string]interface{}, *errors.RestErr)
	GetCurrentUser(int64) (map[string]interface{}, *errors.RestErr)
	GetFriends(int64) (map[string]interface{}, *errors.RestErr)
	MakeFriends(int64, *users.UserFriend) (map[string]interface{}, *errors.RestErr)
	RemoveFriend(int64, int64) (map[string]interface{}, *errors.RestErr)
}

func (s *userService) RemoveFriend(userId int64, friendId int64) (map[string]interface{}, *errors.RestErr) {

	err := users.UserRepository.RemoveFriend(userId, friendId)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *userService) MakeFriends(userId int64, friend *users.UserFriend) (map[string]interface{}, *errors.RestErr) {
	err := users.UserRepository.MakeFriends(userId, friend)
	if err != nil {
		return nil, err
	}

	responseData := response.Data("friend", friend.Id)

	return responseData, nil
}

func (s *userService) GetFriends(userId int64) (map[string]interface{}, *errors.RestErr) {
	friendsList, err := users.UserRepository.GetFriends(userId)
	if err != nil {
		return nil, err
	}

	responseData := response.Data("friends", friendsList)

	return responseData, nil
}

//LoginUser логинит пользователя
func (s *userService) LoginUser(user *users.UserRequest) (map[string]interface{}, *errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	foundUser, err := users.UserRepository.GetUserByEmail(user)
	if err != nil {
		return nil, err
	}

	if err = SecurityService.VerifyPassword(user.Password, foundUser.Password); err != nil {
		return nil, err
	}

	token, err := SecurityService.GenerateToken(foundUser.Id)
	if err != nil {
		return nil, err
	}

	responseUser := &users.ResponseUser{
		FirstName: foundUser.FirstName,
		LastName:  foundUser.LastName,
		Email:     foundUser.Email,
		Interests: foundUser.Interests,
		City:      foundUser.City,
		Age:       foundUser.Age,
		Gender:    foundUser.Gender,
		Token:     token,
	}
	responseData := response.Data("user", responseUser)
	return responseData, nil
}

//RegisterUser регистрирует пользователя
func (s *userService) RegisterUser(user *users.UserRequest) (map[string]interface{}, *errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	hp, err := SecurityService.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hp

	userId, err := users.UserRepository.Save(user)
	if err != nil {
		return nil, err
	}

	token, err := SecurityService.GenerateToken(userId)
	if err != nil {
		return nil, err
	}

	responseUser := &users.ResponseUser{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Interests: user.Interests,
		City:      user.City,
		Age:       user.Age,
		Gender:    user.Gender,
		Token:     token,
	}

	responseData := response.Data("user", responseUser)

	return responseData, nil
}

//GetUsers возвращает пользователей системы
func (s *userService) GetUsers() (map[string]interface{}, *errors.RestErr) {
	usersList, err := users.UserRepository.GetUsers()
	if err != nil {
		return nil, err
	}

	responseData := response.Data("users", usersList)

	return responseData, nil
}

//GetCurrentUser для получения инфы по пользователю
func (s *userService) GetCurrentUser(userId int64) (map[string]interface{}, *errors.RestErr) {
	user, err := users.UserRepository.GetCurrentUser(userId)
	if err != nil {
		return nil, err
	}

	responseData := response.Data("user", user)

	return responseData, nil
}
