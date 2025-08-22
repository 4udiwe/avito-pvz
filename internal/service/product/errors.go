package product

import "errors"

var (
	ErrReceptionAlreadyClosed = errors.New("reception already closed")
	ErrNoPointFound           = errors.New("no point found")
	ErrNoReceptionFound       = errors.New("no reception found")
)
