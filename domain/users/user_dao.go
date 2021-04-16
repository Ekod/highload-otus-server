package users

import (
	"fmt"
	"github.com/Ekod/highload-otus/datasources/mysql/users_db"
	"github.com/Ekod/highload-otus/utils/errors"
	"github.com/Ekod/highload-otus/utils/logger"
)

const (
	queryCreateUser     = "INSERT INTO users(first_name, last_name, age, gender, interests, city, email, password) VALUES(?,?,?,?,?,?,?,?);"
	queryGetUserByEmail = "SELECT id, first_name, last_name, age, gender, interests, city, email, password FROM users WHERE email = ?;"
	queryGetUsers       = "SELECT id, first_name, last_name, age, gender, interests, city, email FROM users;"
	queryGetUserById    = "SELECT first_name, last_name, age, gender, interests, city, email FROM users WHERE id = ?;"
	queryMakeFriends    = "INSERT INTO friends(users_id, friends_id) values(?,?),(?,?);"
	queryGetFriends     = "SELECT users.id as uid, first_name, last_name, age, gender, interests, city, email FROM users JOIN friends ON users.id = friends_id WHERE users_id = ?;"
	queryDeleteFriend   = "DELETE FROM friends WHERE users_id = ? AND friends_id = ?;"
)

var UserRepository userRepositoryInterface = &userRepository{}

type userRepository struct{}

type userRepositoryInterface interface {
	GetFriends(int64) ([]UserFriend, *errors.RestErr)
	GetUserByEmail(*UserRequest) (*User, *errors.RestErr)
	GetCurrentUser(int64) (*ResponseUser, *errors.RestErr)
	GetUsers() ([]ResponseUser, *errors.RestErr)
	Save(*UserRequest) (int64, *errors.RestErr)
	MakeFriends(int64, *UserFriend) *errors.RestErr
	RemoveFriend(int64, int64) *errors.RestErr
}

func (ur *userRepository) RemoveFriend(userId int64, friendId int64) *errors.RestErr {
	stmt, err := users_db.UserClient.Prepare(queryDeleteFriend)
	if err != nil {
		return errors.NewInternalServerError("Server error")
	}
	_, err = stmt.Exec(userId, friendId)
	if err != nil {
		return errors.NewInternalServerError("Server error")
	}

	_, err = stmt.Exec(friendId, userId)
	if err != nil {
		return errors.NewInternalServerError("Server error")
	}

	if err := stmt.Close(); err != nil {
		logger.LogError("Couldn't close the DB statement", err)
		return errors.NewInternalServerError("Server error")
	}

	return nil
}

func (ur *userRepository) MakeFriends(userId int64, friend *UserFriend) *errors.RestErr {
	stmt, err := users_db.UserClient.Prepare(queryMakeFriends)
	if err != nil {
		return errors.NewInternalServerError("Server error")
	}
	_, err = stmt.Exec(userId, friend.Id, friend.Id, userId)
	if err != nil {
		return errors.NewInternalServerError("Server error")
	}
	if err := stmt.Close(); err != nil {
		logger.LogError("Couldn't close the DB statement", err)
		return errors.NewInternalServerError("Server error")
	}
	return nil
}

func (ur *userRepository) GetFriends(userId int64) ([]UserFriend, *errors.RestErr) {
	stmt, err := users_db.UserClient.Prepare(queryGetFriends)
	if err != nil {
		return nil, errors.NewInternalServerError("Server error")
	}

	var friends []UserFriend
	rows, err := stmt.Query(userId)
	if err != nil {
		return nil, errors.NewInternalServerError("Server error")
	}
	for rows.Next() {
		u := UserFriend{}
		err := rows.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Age, &u.Gender, &u.Interests, &u.City, &u.Email)
		if err != nil {
			fmt.Println(err)
			continue
		}
		friends = append(friends, u)
	}

	if err := stmt.Close(); err != nil {
		logger.LogError("Couldn't close the DB statement", err)
		return nil, errors.NewInternalServerError("Server error")
	}
	return friends, nil
}

//GetUserByEmail используется для получения данных о пользователе при логине
func (ur *userRepository) GetUserByEmail(user *UserRequest) (*User, *errors.RestErr) {
	stmt, err := users_db.UserClient.Prepare(queryGetUserByEmail)
	if err != nil {
		return nil, errors.NewInternalServerError("Server error")
	}

	var foundUser User
	if err = stmt.QueryRow(user.Email).Scan(
		&foundUser.Id,
		&foundUser.FirstName,
		&foundUser.LastName,
		&foundUser.Age,
		&foundUser.Gender,
		&foundUser.Interests,
		&foundUser.City,
		&foundUser.Email,
		&foundUser.Password); err != nil {
		return nil, errors.NewInternalServerError("Server error")
	}
	if err := stmt.Close(); err != nil {
		logger.LogError("Couldn't close the DB statement", err)
		return nil, errors.NewInternalServerError("Server error")
	}
	return &foundUser, nil
}

//GetCurrentUser для получения инфы по пользователю
func (ur *userRepository) GetCurrentUser(userId int64) (*ResponseUser, *errors.RestErr) {
	stmt, err := users_db.UserClient.Prepare(queryGetUserById)
	if err != nil {
		return nil, errors.NewInternalServerError("Server error")
	}

	var foundUser ResponseUser
	if err = stmt.QueryRow(userId).Scan(
		&foundUser.FirstName,
		&foundUser.LastName,
		&foundUser.Age,
		&foundUser.Gender,
		&foundUser.Interests,
		&foundUser.City,
		&foundUser.Email); err != nil {
		return nil, errors.NewInternalServerError("Server error")
	}
	if err := stmt.Close(); err != nil {
		logger.LogError("Couldn't close the DB statement", err)
		return nil, errors.NewInternalServerError("Server error")
	}
	return &foundUser, nil
}

//GetUsers для получения всех пользователей
func (ur *userRepository) GetUsers() ([]ResponseUser, *errors.RestErr) {
	stmt, err := users_db.UserClient.Prepare(queryGetUsers)
	if err != nil {
		return nil, errors.NewInternalServerError("Server error")
	}

	var users []ResponseUser
	rows, err := stmt.Query()
	if err != nil {
		return nil, errors.NewInternalServerError("Server error")
	}
	for rows.Next() {
		u := ResponseUser{}
		err := rows.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Age, &u.Gender, &u.Interests, &u.City, &u.Email)
		if err != nil {
			fmt.Println(err)
			continue
		}
		users = append(users, u)
	}

	if err := stmt.Close(); err != nil {
		logger.LogError("Couldn't close the DB statement", err)
		return nil, errors.NewInternalServerError("Server error")
	}

	return users, nil
}

//Save для регистрации пользователя
func (ur *userRepository) Save(user *UserRequest) (int64, *errors.RestErr) {
	stmt, err := users_db.UserClient.Prepare(queryCreateUser)
	if err != nil {
		return 0, errors.NewInternalServerError("Server error")
	}

	result, err := stmt.Exec(user.FirstName, user.LastName, user.Age, user.Gender, user.Interests, user.City, user.Email, user.Password)
	if err != nil {
		return 0, errors.NewInternalServerError("Server error")
	}

	userId, err := result.LastInsertId()
	if err != nil {
		return 0, errors.NewInternalServerError("Server error")
	}
	if err := stmt.Close(); err != nil {
		logger.LogError("Couldn't close the DB statement", err)
		return 0, errors.NewInternalServerError("Server error")
	}
	return userId, nil
}
