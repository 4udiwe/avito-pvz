package post_login

import (
	"context"

	"github.com/4udiwe/avito-pvz/internal/auth"
)

//go:generate go tool mockgen -source=contracts.go -destination=mocks/mock_service.go

type UserService interface {
	Authenticate(ctx context.Context, email, password string) (*auth.Tokens, error)
}
