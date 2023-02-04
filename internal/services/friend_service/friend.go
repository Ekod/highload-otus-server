package friend_service

import (
	"github.com/Ekod/highload-otus/domain/users"
	"github.com/Ekod/highload-otus/utils/errors"
	"github.com/Ekod/highload-otus/utils/response"
)

type FriendRepository interface {
	GetFriends(userId int64) ([]users.UserFriend, *errors.RestErr)
	MakeFriends(int64, *users.UserFriend) *errors.RestErr
	RemoveFriend(int64, int64) *errors.RestErr
}

type Service struct {
	friendRepository FriendRepository
}

func New(friendRepository FriendRepository) *Service {
	return &Service{
		friendRepository: friendRepository,
	}
}

func (s *Service) RemoveFriend(userId int64, friendId int64) (map[string]interface{}, *errors.RestErr) {
	err := s.friendRepository.RemoveFriend(userId, friendId)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *Service) MakeFriends(userId int64, friend *users.UserFriend) (map[string]interface{}, *errors.RestErr) {
	err := s.friendRepository.MakeFriends(userId, friend)
	if err != nil {
		return nil, err
	}

	responseData := response.Data("friend_service", friend.Id)

	return responseData, nil
}

func (s *Service) GetFriends(userId int64) (map[string]interface{}, *errors.RestErr) {
	friendsList, err := s.friendRepository.GetFriends(userId)
	if err != nil {
		return nil, err
	}

	responseData := response.Data("friends", friendsList)

	return responseData, nil
}
