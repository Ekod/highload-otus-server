package repositories

import (
	"database/sql"
	"fmt"
	"github.com/Ekod/highload-otus/domain"
	"github.com/Ekod/highload-otus/domain/users"
	"github.com/Ekod/highload-otus/utils/errors"
	"github.com/Ekod/highload-otus/utils/logger"
)

const (
	querySaveUser       = "INSERT INTO users(first_name, last_name, age, gender, interests, city, email, password) VALUES(?,?,?,?,?,?,?,?);"
	queryGetUserByEmail = "SELECT id, first_name, last_name, age, gender, interests, city, email, password FROM users WHERE email = ?;"
	queryGetUsers       = "SELECT id, first_name, last_name, age, gender, interests, city, email FROM users;"
	queryGetUserById    = "SELECT first_name, last_name, age, gender, interests, city, email FROM users WHERE id = ?;"
	queryLikeSelect     = "SELECT id, first_name, last_name, age, gender, interests, city, email FROM users WHERE first_name LIKE ? AND last_name LIKE ? ORDER BY id;"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) GetUserByEmail(user *users.UserRequest) (*domain.User, *errors.RestErr) {
	stmt, err := ur.db.Prepare(queryGetUserByEmail)
	if err != nil {
		logger.LogError("UserRepository_GetUserByEmail - PrepareQuery", err)
		return nil, errors.NewInternalServerError("Server error")
	}
	defer stmt.Close()

	var foundUser domain.User
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
func (ur *UserRepository) GetCurrentUser(userId int64) (*users.ResponseUser, *errors.RestErr) {
	stmt, err := ur.db.Prepare(queryGetUserById)
	if err != nil {
		logger.LogError("UserRepository_GetCurrentUser - PrepareQuery", err)
		return nil, errors.NewInternalServerError("Server error")
	}
	defer stmt.Close()

	var foundUser users.ResponseUser
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
func (ur *UserRepository) GetUsers() ([]users.ResponseUser, *errors.RestErr) {
	stmt, err := ur.db.Prepare(queryGetUsers)
	if err != nil {
		logger.LogError("UserRepository_GetUsers - PrepareQuery", err)
		return nil, errors.NewInternalServerError("Server error")
	}
	defer stmt.Close()

	var responseUsers []users.ResponseUser

	rows, err := stmt.Query()
	if err != nil {
		logger.LogError("UserRepository_GetUsers - Query", err)
		return nil, errors.NewInternalServerError("Server error")
	}
	defer rows.Close()

	for rows.Next() {
		user := users.ResponseUser{}
		err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Age, &user.Gender, &user.Interests, &user.City, &user.Email)
		if err != nil {
			logger.LogError("UserRepository_GetUsers - RowsNext", err)
			continue
		}

		responseUsers = append(responseUsers, user)
	}

	return responseUsers, nil
}
func (ur *UserRepository) SaveUser(user *users.UserRequest) (int64, *errors.RestErr) {
	stmt, err := ur.db.Prepare(querySaveUser)
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
func (ur *UserRepository) GetUsersByFullName(firstName, lastName string) ([]users.ResponseUser, *errors.RestErr) {
	stmt, err := ur.db.Prepare(queryLikeSelect)
	if err != nil {
		logger.LogError("UserRepository_GetUsersByFullName - PrepareQuery", err)
		return nil, errors.NewInternalServerError("Server error")
	}
	defer stmt.Close()

	var responseUsers []users.ResponseUser

	rows, err := stmt.Query(fmt.Sprintf("%s%", firstName), fmt.Sprintf("%s%", lastName))
	if err != nil {
		logger.LogError("UserRepository_GetUsersByFullName - Query", err)
		return nil, errors.NewInternalServerError("Server error")
	}
	defer rows.Close()

	for rows.Next() {
		user := users.ResponseUser{}
		err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Age, &user.Gender, &user.Interests, &user.City, &user.Email)
		if err != nil {
			logger.LogError("UserRepository_GetUsersByFullName - RowsNext", err)
			continue
		}

		responseUsers = append(responseUsers, user)
	}

	return responseUsers, nil
}
