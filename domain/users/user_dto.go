package users

type UserFriend struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	Interests string `json:"interests"`
	City      string `json:"city"`
	Email     string `json:"email"`
}

type UserRequest struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Age       int    `json:"age" validate:"required,number"`
	Gender    string `json:"gender" validate:"required,oneof=male female"`
	Interests string `json:"interests" validate:"required"`
	City      string `json:"city" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
}

type UserResponse struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
	Gender    string `json:"gender"`
	Interests string `json:"interests"`
	City      string `json:"city"`
	Email     string `json:"email"`
	Token     string `json:"token,omitempty"`
}
