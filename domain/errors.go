package domain

import "errors"

var (
	ErrNotFound      = errors.New("Your requested Item is not found")
	ErrConflict      = errors.New("Your Item already exist")
	ErrBadParamInput = errors.New("Given Param is not valid")
)
