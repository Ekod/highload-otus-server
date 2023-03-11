package errors

import (
	"errors"
	"fmt"
	"net/http"
)

func (r *RestErr) Error() string {
	return fmt.Sprintf("%d %s %s", r.Status, r.Message, r.Err)

}

type RestErr struct {
	Message      string `json:"message"`
	Status       int    `json:"status"`
	Err          string `json:"error"`
	DebugMessage string `json:"-"`
}

func NewBadRequestError(message, debugMessage string) *RestErr {
	return &RestErr{
		Message:      message,
		Status:       http.StatusBadRequest,
		Err:          "bad_request",
		DebugMessage: debugMessage,
	}
}

func NewHandlerBadRequestError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusBadRequest,
		Err:     "bad_request",
	}
}

func NewNotFoundError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusNotFound,
		Err:     "not_found",
	}
}

func NewInternalServerError(message, debugMessage string) *RestErr {
	return &RestErr{
		Message:      message,
		Status:       http.StatusInternalServerError,
		Err:          "internal_server_error",
		DebugMessage: debugMessage,
	}
}

func NewUnauthorizedError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusUnauthorized,
		Err:     "unauthorized",
	}
}

func ParseError(err error) RestErr {
	var restErr *RestErr

	errors.As(err, &restErr)

	return *restErr
}
