package users

import (
	"github.com/Ekod/highload-otus/utils/errors"
	"strings"
)

type UserFriend struct {
	Id        int64  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	Interests string `json:"interests"`
	City      string `json:"city"`
	Email     string `json:"email"`
}

type SearchLikeUser struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

//UserRequest используется для регистрации пользователя
type UserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	Interests string `json:"interests"`
	City      string `json:"city"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

//User используется для отдачи данных о пользователе
type User struct {
	Id        int64  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	Interests string `json:"interests"`
	City      string `json:"city"`
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"`
}

type ResponseUser struct {
	Id        int64  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	Interests string `json:"interests"`
	City      string `json:"city"`
	Email     string `json:"email"`
	Token     string `json:"token,omitempty"`
}

//SecurityUser используется при логине пользователя
type SecurityUser struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

//Validate валидирует поля приходящего запроса
func (user *UserRequest) Validate() *errors.RestErr {
	trimmedEmail := strings.TrimSpace(user.Email)
	if trimmedEmail == "" {
		return errors.NewBadRequestError("Email or password is invalid")
	}

	trimmedPassword := strings.TrimSpace(user.Password)
	if trimmedPassword == "" {
		return errors.NewBadRequestError("Email or password is invalid")
	}

	return nil
}
