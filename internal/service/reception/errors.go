package reception

import "errors"

var (
	ErrLastReceptionNotClosed    = errors.New("last reception not closed")
	ErrNoPointFound              = errors.New("no point found")
	ErrCannotCloseEmptyReception = errors.New("cannot close empty reception")
	ErrNoReceptionFound          = errors.New("no reception found")
)
