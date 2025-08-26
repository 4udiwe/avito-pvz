package user

import "errors"

var (
	ErrNoUserFound         = errors.New("no user found")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrInvalidCredentials  = errors.New("invalid credentials")
)
