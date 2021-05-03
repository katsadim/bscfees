package handler

import (
	"fmt"
	"net/http"
	"time"
)

type Error struct {
	Message string `json:"message"`
	Time    string `json:"time"`
	Code    int    `json:"-"`
}

func NewError(m string, c int) *Error {
	return &Error{
		Message: m,
		Time:    time.Now().UTC().Format(time.RFC3339),
		Code:    c,
	}
}

func NewBadRequestError(m string) *Error {
	return NewError(
		m,
		http.StatusBadRequest,
	)
}

func NewEmptyQueryParameterError(param string) *Error {
	return NewBadRequestError(fmt.Sprintf("[%s] is empty", param))
}

func NewInternalError(m string) *Error {
	return NewError(m, http.StatusInternalServerError)
}
