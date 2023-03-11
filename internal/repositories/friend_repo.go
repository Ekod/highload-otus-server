package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Ekod/highload-otus/domain/users"
	"github.com/Ekod/highload-otus/utils/errors"
)

const (
	queryMakeFriends  = "INSERT INTO friends(user_id, friend_id) values(?,?),(?,?);"
	queryGetFriends   = "SELECT users.id as uid, first_name, last_name, age, gender, interests, city, email FROM users JOIN friends ON users.id = friends_id WHERE users_id = ?;"
	queryDeleteFriend = "DELETE FROM friends WHERE user_id = ? AND friend_id = ?;"
)

type FriendRepository struct {
	db *sql.DB
}

func NewFriendRepository(db *sql.DB) *FriendRepository {
	return &FriendRepository{db: db}
}

func (us *FriendRepository) GetFriends(ctx context.Context, userId int) ([]users.UserFriend, error) {
	stmt, err := us.db.Prepare(queryGetFriends)
	if err != nil {
		debugMessageError := fmt.Sprintf("[ERROR] UserRepository_GetFriends - PrepareQuery: %s", err)

		return nil, errors.NewInternalServerError("Server error", debugMessageError)
	}
	defer stmt.Close()

	var friends []users.UserFriend

	rows, err := stmt.QueryContext(ctx, userId)
	if err != nil {
		debugMessageError := fmt.Sprintf("[ERROR] UserRepository_GetFriends - Query: %s", err)

		return nil, errors.NewInternalServerError("Server error", debugMessageError)
	}
	defer rows.Close()

	for rows.Next() {
		var u users.UserFriend

		err = rows.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Age, &u.Gender, &u.Interests, &u.City, &u.Email)
		if err != nil {
			continue
		}

		friends = append(friends, u)
	}

	return friends, nil
}

func (us *FriendRepository) MakeFriends(ctx context.Context, userId int, friendID int) (returnErr error) {
	tx, err := us.db.Begin()
	if err != nil {
		debugMessageError := fmt.Sprintf("[ERROR] UserRepository_MakeFriends - Transaction Begin: %s", err)
		returnErr = errors.NewInternalServerError("Server error", debugMessageError)

		return returnErr
	}
	defer func() {
		if returnErr != nil {
			tx.Rollback()
		}

		tx.Commit()
	}()

	stmt, err := tx.Prepare(queryMakeFriends)
	if err != nil {
		debugMessageError := fmt.Sprintf("[ERROR] UserRepository_MakeFriends - PrepareQuery: %s", err)
		returnErr = errors.NewInternalServerError("Server error", debugMessageError)

		return returnErr
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, userId, friendID, friendID, userId)
	if err != nil {
		debugMessageError := fmt.Sprintf("[ERROR] UserRepository_MakeFriends - ExecQuery: %s", err)
		returnErr = errors.NewInternalServerError("Server error", debugMessageError)

		return returnErr
	}

	return nil
}

func (us *FriendRepository) RemoveFriend(ctx context.Context, userId int, friendId int) (returnErr error) {
	tx, err := us.db.Begin()
	if err != nil {
		debugMessageError := fmt.Sprintf("[ERROR] UserRepository_RemoveFriend - Transaction Begin: %s", err)
		returnErr = errors.NewInternalServerError("Server error", debugMessageError)

		return returnErr
	}
	defer func() {
		if returnErr != nil {
			tx.Rollback()
		}

		tx.Commit()
	}()

	stmt, err := tx.Prepare(queryDeleteFriend)
	if err != nil {
		debugMessageError := fmt.Sprintf("[ERROR] UserRepository_RemoveFriend - PrepareQuery: %s", err)
		returnErr = errors.NewInternalServerError("Server error", debugMessageError)

		return returnErr
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, userId, friendId)
	if err != nil {
		debugMessageError := fmt.Sprintf("[ERROR] UserRepository_RemoveFriend - ExecQuery1: %s", err)
		returnErr = errors.NewInternalServerError("Server error", debugMessageError)

		return returnErr
	}

	_, err = stmt.Exec(friendId, userId)
	if err != nil {
		debugMessageError := fmt.Sprintf("[ERROR] UserRepository_RemoveFriend - ExecQuery2: %s", err)
		returnErr = errors.NewInternalServerError("Server error", debugMessageError)

		return returnErr
	}

	return nil
}
