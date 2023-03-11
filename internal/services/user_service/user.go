package user_service

import (
	"context"
	"github.com/Ekod/highload-otus/domain"
	"github.com/Ekod/highload-otus/domain/users"
	"github.com/Ekod/highload-otus/utils/security"
)

type UserRepository interface {
	GetUserByEmail(context.Context, *users.UserRequest) (*domain.User, error)
	GetCurrentUser(context.Context, int) (*users.UserResponse, error)
	GetUsers(context.Context) ([]users.UserResponse, error)
	SaveUser(context.Context, *users.UserRequest) (int, error)
	GetUsersByFullName(context.Context, string, string) ([]users.UserResponse, error)
}

type Service struct {
	userRepository UserRepository
}

func New(userRepository UserRepository) *Service {
	return &Service{
		userRepository: userRepository,
	}
}

func (s *Service) LoginUser(ctx context.Context, user *users.UserRequest) (*users.UserResponse, error) {
	foundUser, err := s.userRepository.GetUserByEmail(ctx, user)
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

	responseUser := users.UserResponse{
		FirstName: foundUser.FirstName,
		LastName:  foundUser.LastName,
		Email:     foundUser.Email,
		Interests: foundUser.Interests,
		City:      foundUser.City,
		Age:       foundUser.Age,
		Gender:    foundUser.Gender,
		Token:     token,
	}

	return &responseUser, nil
}

func (s *Service) RegisterUser(ctx context.Context, user *users.UserRequest) (*users.UserResponse, error) {
	hp, err := security.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hp

	userId, err := s.userRepository.SaveUser(ctx, user)
	if err != nil {
		return nil, err
	}

	token, err := security.GenerateToken(userId)
	if err != nil {
		return nil, err
	}

	responseUser := users.UserResponse{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Interests: user.Interests,
		City:      user.City,
		Age:       user.Age,
		Gender:    user.Gender,
		Token:     token,
	}

	return &responseUser, nil
}

func (s *Service) GetUsers(ctx context.Context) ([]users.UserResponse, error) {
	usersList, err := s.userRepository.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	return usersList, nil
}

func (s *Service) GetCurrentUser(ctx context.Context, userId int) (*users.UserResponse, error) {
	user, err := s.userRepository.GetCurrentUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetUsersByFullName(ctx context.Context, firstName, lastName string) ([]users.UserResponse, error) {
	friendsList, err := s.userRepository.GetUsersByFullName(ctx, firstName, lastName)
	if err != nil {
		return nil, err
	}

	return friendsList, nil
}
