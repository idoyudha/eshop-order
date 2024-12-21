package v1

import (
	"net/http"
)

type restError struct {
	Code  int          `json:"code"`
	Error errorMessage `json:"error"`
}

type errorMessage struct {
	Message string `json:"message"`
	Causes  error  `json:"causes"`
}

func newBadRequestError(message string) *restError {
	return &restError{
		Code: http.StatusBadRequest,
		Error: errorMessage{
			Message: message,
		},
	}
}

func newNotFoundError(message string) *restError {
	return &restError{
		Code: http.StatusNotFound,
		Error: errorMessage{
			Message: message,
		},
	}
}

func newInternalServerError(message string) *restError {
	return &restError{
		Code: http.StatusInternalServerError,
		Error: errorMessage{
			Message: message,
		},
	}
}

func newUnauthorizedError(message string) *restError {
	return &restError{
		Code: http.StatusUnauthorized,
		Error: errorMessage{
			Message: message,
		},
	}
}
