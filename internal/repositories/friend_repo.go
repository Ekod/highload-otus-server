package repositories

import (
	"database/sql"
	"github.com/Ekod/highload-otus/domain/users"
	"github.com/Ekod/highload-otus/utils/errors"
	"github.com/Ekod/highload-otus/utils/logger"
)

const (
	queryMakeFriends  = "INSERT INTO friends(users_id, friends_id) values(?,?),(?,?);"
	queryGetFriends   = "SELECT users.id as uid, first_name, last_name, age, gender, interests, city, email FROM users JOIN friends ON users.id = friends_id WHERE users_id = ?;"
	queryDeleteFriend = "DELETE FROM friends WHERE users_id = ? AND friends_id = ?;"
)

type FriendRepository struct {
	db *sql.DB
}

func NewFriendRepository(db *sql.DB) *FriendRepository {
	return &FriendRepository{db: db}
}

func (us *FriendRepository) GetFriends(userId int64) ([]users.UserFriend, *errors.RestErr) {
	stmt, err := us.db.Prepare(queryGetFriends)
	if err != nil {
		logger.LogError("UserRepository_GetFriends - PrepareQuery", err)
		return nil, errors.NewInternalServerError("Server error")
	}
	defer stmt.Close()

	var friends []users.UserFriend

	rows, err := stmt.Query(userId)
	if err != nil {
		logger.LogError("UserRepository_GetFriends - Query", err)
		return nil, errors.NewInternalServerError("Server error")
	}
	defer rows.Close()

	for rows.Next() {
		u := users.UserFriend{}
		err := rows.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Age, &u.Gender, &u.Interests, &u.City, &u.Email)
		if err != nil {
			logger.LogError("UserRepository_GetFriends - RowsNext", err)
			continue
		}
		friends = append(friends, u)
	}

	return friends, nil
}

func (us *FriendRepository) MakeFriends(userId int64, friend *users.UserFriend) (returnErr *errors.RestErr) {
	tx, err := us.db.Begin()
	if err != nil {
		logger.LogError("UserRepository_MakeFriends - Transaction Begin", err)
		returnErr = errors.NewInternalServerError("Server error")
		return returnErr
	}
	defer func() {
		if returnErr != nil {
			if err = tx.Rollback(); err != nil {
				logger.LogError("UserRepository_MakeFriends - Transaction Rollback", err)
			}
		}
		if err = tx.Commit(); err != nil {
			logger.LogError("UserRepository_MakeFriends - Transaction Commit", err)
		}
	}()

	stmt, err := tx.Prepare(queryMakeFriends)
	if err != nil {
		logger.LogError("UserRepository_MakeFriends - PrepareQuery", err)
		returnErr = errors.NewInternalServerError("Server error")
		return returnErr
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId, friend.Id, friend.Id, userId)
	if err != nil {
		logger.LogError("UserRepository_MakeFriends - ExecQuery", err)
		returnErr = errors.NewInternalServerError("Server error")
		return returnErr
	}

	return nil
}

func (us *FriendRepository) RemoveFriend(userId int64, friendId int64) (returnErr *errors.RestErr) {
	tx, err := us.db.Begin()
	if err != nil {
		logger.LogError("UserRepository_RemoveFriend - Transaction Begin", err)
		returnErr = errors.NewInternalServerError("Server error")
		return returnErr
	}
	defer func() {
		if returnErr != nil {
			if err = tx.Rollback(); err != nil {
				logger.LogError("UserRepository_RemoveFriend - Transaction Rollback", err)
			}
		}
		if err = tx.Commit(); err != nil {
			logger.LogError("UserRepository_RemoveFriend - Transaction Commit", err)
		}
	}()

	stmt, err := tx.Prepare(queryDeleteFriend)
	if err != nil {
		logger.LogError("UserRepository_RemoveFriend - PrepareQuery", err)
		returnErr = errors.NewInternalServerError("Server error")
		return returnErr
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId, friendId)
	if err != nil {
		logger.LogError("UserRepository_RemoveFriend - ExecQuery1", err)
		returnErr = errors.NewInternalServerError("Server error")
		return returnErr
	}

	_, err = stmt.Exec(friendId, userId)
	if err != nil {
		logger.LogError("UserRepository_RemoveFriend - ExecQuery2", err)
		returnErr = errors.NewInternalServerError("Server error")
		return returnErr
	}

	return nil
}
