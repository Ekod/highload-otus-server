package user_service

import (
	"github.com/Ekod/highload-otus/domain"
	"github.com/Ekod/highload-otus/domain/users"
	"github.com/Ekod/highload-otus/utils/errors"
	"github.com/Ekod/highload-otus/utils/response"
	"github.com/Ekod/highload-otus/utils/security"
)

type UserRepository interface {
	GetUserByEmail(*users.UserRequest) (*domain.User, *errors.RestErr)
	GetCurrentUser(int64) (*users.ResponseUser, *errors.RestErr)
	GetUsers() ([]users.ResponseUser, *errors.RestErr)
	SaveUser(*users.UserRequest) (int64, *errors.RestErr)
	GetUsersByFullName(string, string) ([]users.ResponseUser, *errors.RestErr)
}

type Service struct {
	userRepository UserRepository
}

func New(userRepository UserRepository) *Service {
	return &Service{
		userRepository: userRepository,
	}
}

func (s *Service) LoginUser(user *users.UserRequest) (map[string]interface{}, *errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	foundUser, err := s.userRepository.GetUserByEmail(user)
	if err != nil {
		return nil, err
	}

	if err = security.VerifyPassword(user.Password, foundUser.Password); err != nil {
		return nil, err
	}

	token, err := security.GenerateToken(foundUser.Id)
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
	responseData := response.Data("user_service", responseUser)
	return responseData, nil
}

func (s *Service) RegisterUser(user *users.UserRequest) (map[string]interface{}, *errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	hp, err := security.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hp

	userId, err := s.userRepository.SaveUser(user)
	if err != nil {
		return nil, err
	}

	token, err := security.GenerateToken(userId)
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

	responseData := response.Data("user_service", responseUser)

	return responseData, nil
}

func (s *Service) GetUsers() (map[string]interface{}, *errors.RestErr) {
	usersList, err := s.userRepository.GetUsers()
	if err != nil {
		return nil, err
	}

	responseData := response.Data("users", usersList)

	return responseData, nil
}

func (s *Service) GetCurrentUser(userId int64) (map[string]interface{}, *errors.RestErr) {
	user, err := s.userRepository.GetCurrentUser(userId)
	if err != nil {
		return nil, err
	}

	responseData := response.Data("user_service", user)

	return responseData, nil
}

func (s *Service) GetUsersByFullName(firstName, lastName string) (map[string][]users.ResponseUser, *errors.RestErr) {
	friendsList, err := s.userRepository.GetUsersByFullName(firstName, lastName)
	if err != nil {
		return nil, err
	}

	responseData := make(map[string][]users.ResponseUser)
	responseData["users"] = friendsList

	return responseData, nil
}
