package post_login

import (
	"context"

	"github.com/4udiwe/avito-pvz/internal/auth"
)

type UserService interface {
	Authenticate(ctx context.Context, email, password string) (*auth.Tokens, error)
}
