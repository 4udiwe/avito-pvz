package post_register

import (
	"context"

	"github.com/4udiwe/avito-pvz/internal/auth"
	"github.com/4udiwe/avito-pvz/internal/entity"
)

//go:generate go tool mockgen -source=contracts.go -destination=mocks/mock_service.go

type UserService interface {
	Register(ctx context.Context, email, password string, role entity.UserRole) (*auth.Tokens, error)
}
