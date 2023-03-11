package friend_service

import (
	"context"
	"github.com/Ekod/highload-otus/domain/users"
)

type FriendRepository interface {
	GetFriends(context.Context, int) ([]users.UserFriend, error)
	MakeFriends(context.Context, int, int) error
	RemoveFriend(context.Context, int, int) error
}

type Service struct {
	friendRepository FriendRepository
}

func New(friendRepository FriendRepository) *Service {
	return &Service{
		friendRepository: friendRepository,
	}
}

func (s *Service) RemoveFriend(ctx context.Context, userId int, friendId int) error {
	err := s.friendRepository.RemoveFriend(ctx, userId, friendId)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) MakeFriends(ctx context.Context, userId int, friendID int) (int, error) {
	err := s.friendRepository.MakeFriends(ctx, userId, friendID)
	if err != nil {
		return 0, err
	}

	return friendID, nil
}

func (s *Service) GetFriends(ctx context.Context, userId int) ([]users.UserFriend, error) {
	friendsList, err := s.friendRepository.GetFriends(ctx, userId)
	if err != nil {
		return nil, err
	}

	return friendsList, nil
}
