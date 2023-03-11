package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Ekod/highload-otus/domain"
	"github.com/Ekod/highload-otus/domain/users"
	"github.com/Ekod/highload-otus/utils/errors"
)

const (
	querySaveUser       = "INSERT INTO users(first_name, last_name, age, gender, interests, city, email, password, created_at, updated_at) VALUES(?,?,?,?,?,?,?,?, NOW(), NOW())"
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

func (ur *UserRepository) GetUserByEmail(ctx context.Context, user *users.UserRequest) (*domain.User, error) {
	stmt, err := ur.db.Prepare(queryGetUserByEmail)
	if err != nil {
		debugMessageError := fmt.Sprintf("[ERROR] UserRepository_GetUserByEmail - PrepareQuery: %s", err)

		return nil, errors.NewInternalServerError("Server error", debugMessageError)
	}
	defer stmt.Close()

	var foundUser domain.User

	err = stmt.QueryRowContext(ctx, user.Email).Scan(
		&foundUser.Id,
		&foundUser.FirstName,
		&foundUser.LastName,
		&foundUser.Age,
		&foundUser.Gender,
		&foundUser.Interests,
		&foundUser.City,
		&foundUser.Email,
		&foundUser.Password)
	if err != nil {
		debugMessageError := fmt.Sprintf("[ERROR] UserRepository_GetUserByEmail - Scan: %s", err)

		return nil, errors.NewInternalServerError("Server error", debugMessageError)
	}

	return &foundUser, nil
}
func (ur *UserRepository) GetCurrentUser(ctx context.Context, userId int) (*users.UserResponse, error) {
	stmt, err := ur.db.Prepare(queryGetUserById)
	if err != nil {
		debugMessageError := fmt.Sprintf("[ERROR] UserRepository_GetCurrentUser - PrepareQuery: %s", err)

		return nil, errors.NewInternalServerError("Server error", debugMessageError)
	}
	defer stmt.Close()

	var foundUser users.UserResponse

	err = stmt.QueryRowContext(ctx, userId).Scan(
		&foundUser.FirstName,
		&foundUser.LastName,
		&foundUser.Age,
		&foundUser.Gender,
		&foundUser.Interests,
		&foundUser.City,
		&foundUser.Email)
	if err != nil {
		debugMessageError := fmt.Sprintf("[ERROR] UserRepository_GetCurrentUser - Scan: %s", err)

		return nil, errors.NewInternalServerError("Server error", debugMessageError)
	}

	return &foundUser, nil
}
func (ur *UserRepository) GetUsers(ctx context.Context) ([]users.UserResponse, error) {
	stmt, err := ur.db.Prepare(queryGetUsers)
	if err != nil {
		debugMessageError := fmt.Sprintf("[ERROR] UserRepository_GetUsers - PrepareQuery: %s", err)

		return nil, errors.NewInternalServerError("Server error", debugMessageError)
	}
	defer stmt.Close()

	var responseUsers []users.UserResponse

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		debugMessageError := fmt.Sprintf("[ERROR] UserRepository_GetUsers - Query: %s", err)

		return nil, errors.NewInternalServerError("Server error", debugMessageError)
	}
	defer rows.Close()

	for rows.Next() {
		var user users.UserResponse

		err = rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Age, &user.Gender, &user.Interests, &user.City, &user.Email)
		if err != nil {
			continue
		}

		responseUsers = append(responseUsers, user)
	}

	return responseUsers, nil
}
func (ur *UserRepository) SaveUser(ctx context.Context, user *users.UserRequest) (int, error) {
	stmt, err := ur.db.Prepare(querySaveUser)
	if err != nil {
		debugMessageError := fmt.Sprintf("[ERROR] UserRepository_Save - PrepareQuery: %s", err)

		return 0, errors.NewInternalServerError("Server error", debugMessageError)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx,
		user.FirstName,
		user.LastName,
		user.Age,
		user.Gender,
		user.Interests,
		user.City,
		user.Email,
		user.Password,
	)
	if err != nil {
		debugMessageError := fmt.Sprintf("[ERROR] UserRepository_Save - ExecQuery: %s", err)

		return 0, errors.NewInternalServerError("Server error", debugMessageError)
	}

	userId, err := result.LastInsertId()
	if err != nil {
		debugMessageError := fmt.Sprintf("[ERROR] UserRepository_Save - LastInsertId: %s", err)

		return 0, errors.NewInternalServerError("Server error", debugMessageError)
	}

	return int(userId), nil
}
func (ur *UserRepository) GetUsersByFullName(ctx context.Context, firstName, lastName string) ([]users.UserResponse, error) {
	stmt, err := ur.db.Prepare(queryLikeSelect)
	if err != nil {
		debugMessageError := fmt.Sprintf("[ERROR] UserRepository_GetUsersByFullName - PrepareQuery: %s", err)

		return nil, errors.NewInternalServerError("Server error", debugMessageError)
	}
	defer stmt.Close()

	var responseUsers []users.UserResponse

	rows, err := stmt.QueryContext(ctx, fmt.Sprintf("%s%", firstName), fmt.Sprintf("%s%", lastName))
	if err != nil {
		debugMessageError := fmt.Sprintf("[ERROR] UserRepository_GetUsersByFullName - Query: %s", err)

		return nil, errors.NewInternalServerError("Server error", debugMessageError)
	}
	defer rows.Close()

	for rows.Next() {
		var user users.UserResponse

		err = rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Age, &user.Gender, &user.Interests, &user.City, &user.Email)
		if err != nil {
			continue
		}

		responseUsers = append(responseUsers, user)
	}

	return responseUsers, nil
}
