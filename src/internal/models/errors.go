package models

import (
	"errors"
)

func NewUserErr(text string) error {
	return &UserError{s: text, inner: nil}
}

type UserError struct {
	inner error
	s     string
}

func (e UserError) Error() string {
	return e.s
}

func (e *UserError) Unwrap() error {
	return e.inner
}

var (
	ErrNotFound            = NewUserErr("item not found")
	ErrInternalServ        = NewUserErr("internal server error")
	ErrDecodingRequest     = NewUserErr("broken request")
	ErrDuplicateuserData   = NewUserErr("user with this login already exists")
	ErrDuplicateMarkupType = NewUserErr("ID of this markup already exists")
	ErrViolatingKeyAnnot   = NewUserErr("there is no annot type for this anotattion")
)

func GetUserError(err error) error { // error which will be returned to user
	var userError *UserError

	if errors.As(err, &userError) {
		return *userError
	} else {
		return ErrInternalServ
	}
}
