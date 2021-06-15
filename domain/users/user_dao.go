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
	queryLikeSelect     = "SELECT id, first_name, last_name, age, gender, interests, city, email FROM users WHERE first_name LIKE ? AND last_name LIKE ? ORDER BY id;"
)

var UserRepository userRepositoryInterface = &userRepository{}

type userRepository struct{}

type userRepositoryInterface interface {
	GetFriends(int64) ([]UserFriend, *errors.RestErr)
	GetUserByEmail(*UserRequest) (*User, *errors.RestErr)
	GetCurrentUser(int64) (*ResponseUser, *errors.RestErr)
	GetUsers() ([]ResponseUser, *errors.RestErr)
	SaveUser(*UserRequest) (int64, *errors.RestErr)
	MakeFriends(int64, *UserFriend) *errors.RestErr
	RemoveFriend(int64, int64) *errors.RestErr
	GetUsersByFullName(string, string) ([]ResponseUser, *errors.RestErr)
}

func (ur *userRepository) GetUsersByFullName(firstName, lastName string) ([]ResponseUser, *errors.RestErr) {
	stmt, err := users_db.UserClientSlave1.Prepare(queryLikeSelect)
	if err != nil {
		logger.LogError("UserRepository_GetUsersByFullName - PrepareQuery", err)
		return nil, errors.NewInternalServerError("Server error")
	}
	defer stmt.Close()

	var users []ResponseUser

	rows, err := stmt.Query(fmt.Sprintf("%s%", firstName), fmt.Sprintf("%s%", lastName))
	if err != nil {
		logger.LogError("UserRepository_GetUsersByFullName - Query", err)
		return nil, errors.NewInternalServerError("Server error")
	}
	defer rows.Close()

	for rows.Next() {
		u := ResponseUser{}
		err := rows.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Age, &u.Gender, &u.Interests, &u.City, &u.Email)
		if err != nil {
			logger.LogError("UserRepository_GetUsersByFullName - RowsNext", err)
			continue
		}
		users = append(users, u)
	}

	return users, nil
}

func (ur *userRepository) RemoveFriend(userId int64, friendId int64) *errors.RestErr {
	stmt, err := users_db.UserClientMaster.Prepare(queryDeleteFriend)
	if err != nil {
		logger.LogError("UserRepository_RemoveFriend - PrepareQuery", err)
		return errors.NewInternalServerError("Server error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId, friendId)
	if err != nil {
		logger.LogError("UserRepository_RemoveFriend - ExecQuery1", err)
		return errors.NewInternalServerError("Server error")
	}

	_, err = stmt.Exec(friendId, userId)
	if err != nil {
		logger.LogError("UserRepository_RemoveFriend - ExecQuery2", err)
		return errors.NewInternalServerError("Server error")
	}

	return nil
}

func (ur *userRepository) MakeFriends(userId int64, friend *UserFriend) *errors.RestErr {
	stmt, err := users_db.UserClientMaster.Prepare(queryMakeFriends)
	if err != nil {
		logger.LogError("UserRepository_MakeFriends - PrepareQuery", err)
		return errors.NewInternalServerError("Server error")
	}
	defer stmt.Close()

	_, err = stmt.Exec(userId, friend.Id, friend.Id, userId)
	if err != nil {
		logger.LogError("UserRepository_MakeFriends - ExecQuery", err)
		return errors.NewInternalServerError("Server error")
	}

	return nil
}

func (ur *userRepository) GetFriends(userId int64) ([]UserFriend, *errors.RestErr) {
	stmt, err := users_db.UserClientMaster.Prepare(queryGetFriends)
	if err != nil {
		logger.LogError("UserRepository_GetFriends - PrepareQuery", err)
		return nil, errors.NewInternalServerError("Server error")
	}
	defer stmt.Close()

	var friends []UserFriend

	rows, err := stmt.Query(userId)
	if err != nil {
		logger.LogError("UserRepository_GetFriends - Query", err)
		return nil, errors.NewInternalServerError("Server error")
	}
	defer rows.Close()

	for rows.Next() {
		u := UserFriend{}
		err := rows.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Age, &u.Gender, &u.Interests, &u.City, &u.Email)
		if err != nil {
			logger.LogError("UserRepository_GetFriends - RowsNext", err)
			continue
		}
		friends = append(friends, u)
	}

	return friends, nil
}

//GetUserByEmail используется для получения данных о пользователе при логине
func (ur *userRepository) GetUserByEmail(user *UserRequest) (*User, *errors.RestErr) {
	stmt, err := users_db.UserClientMaster.Prepare(queryGetUserByEmail)
	if err != nil {
		logger.LogError("UserRepository_GetUserByEmail - PrepareQuery", err)
		return nil, errors.NewInternalServerError("Server error")
	}
	defer stmt.Close()

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
		logger.LogError("UserRepository_GetUserByEmail - Scan", err)
		return nil, errors.NewInternalServerError("Server error")
	}

	return &foundUser, nil
}

//GetCurrentUser для получения инфы по пользователю
func (ur *userRepository) GetCurrentUser(userId int64) (*ResponseUser, *errors.RestErr) {
	stmt, err := users_db.UserClientMaster.Prepare(queryGetUserById)
	if err != nil {
		logger.LogError("UserRepository_GetCurrentUser - PrepareQuery", err)
		return nil, errors.NewInternalServerError("Server error")
	}
	defer stmt.Close()

	var foundUser ResponseUser
	if err = stmt.QueryRow(userId).Scan(
		&foundUser.FirstName,
		&foundUser.LastName,
		&foundUser.Age,
		&foundUser.Gender,
		&foundUser.Interests,
		&foundUser.City,
		&foundUser.Email); err != nil {
		logger.LogError("UserRepository_GetCurrentUser - Scan", err)
		return nil, errors.NewInternalServerError("Server error")
	}

	return &foundUser, nil
}

//GetUsers для получения всех пользователей
func (ur *userRepository) GetUsers() ([]ResponseUser, *errors.RestErr) {
	stmt, err := users_db.UserClientSlave1.Prepare(queryGetUsers)
	if err != nil {
		logger.LogError("UserRepository_GetUsers - PrepareQuery", err)
		return nil, errors.NewInternalServerError("Server error")
	}
	defer stmt.Close()

	var users []ResponseUser

	rows, err := stmt.Query()
	if err != nil {
		logger.LogError("UserRepository_GetUsers - Query", err)
		return nil, errors.NewInternalServerError("Server error")
	}
	defer rows.Close()

	for rows.Next() {
		u := ResponseUser{}
		err := rows.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Age, &u.Gender, &u.Interests, &u.City, &u.Email)
		if err != nil {
			logger.LogError("UserRepository_GetUsers - RowsNext", err)
			continue
		}
		users = append(users, u)
	}

	return users, nil
}

//SaveUser для регистрации пользователя
func (ur *userRepository) SaveUser(user *UserRequest) (int64, *errors.RestErr) {
	stmt, err := users_db.UserClientMaster.Prepare(queryCreateUser)
	if err != nil {
		logger.LogError("UserRepository_Save - PrepareQuery", err)
		return 0, errors.NewInternalServerError("Server error")
	}
	defer stmt.Close()

	result, err := stmt.Exec(user.FirstName, user.LastName, user.Age, user.Gender, user.Interests, user.City, user.Email, user.Password)
	if err != nil {
		logger.LogError("UserRepository_Save - ExecQuery", err)
		return 0, errors.NewInternalServerError("Server error")
	}

	userId, err := result.LastInsertId()
	if err != nil {
		logger.LogError("UserRepository_Save - LastInsertId", err)
		return 0, errors.NewInternalServerError("Server error")
	}

	return userId, nil
}
