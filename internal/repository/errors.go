package repository

import "errors"

var (
	ErrNoCityFound  = errors.New("no city found")
	ErrNoPointFound = errors.New("no point found")

	ErrUserAlreadyExists = errors.New("user already exists")
	ErrNoUserFound       = errors.New("no user found")

	ErrLastReceptionNotClosed    = errors.New("last reception not closed")
	ErrCannotCloseEmptyReception = errors.New("cannot close empty reception")
	ErrNoReceptionFound          = errors.New("no reception found")

	ErrNoProductFound = errors.New("no product found")
)
